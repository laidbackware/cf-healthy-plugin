package sig_check

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/plugin"
	logcache "code.cloudfoundry.org/go-log-cache/v3"
	logcache_v1 "code.cloudfoundry.org/go-log-cache/v3/rpc/logcache_v1"
	"code.cloudfoundry.org/go-loggregator/v10/rpc/loggregator_v2"
	logHttp "github.com/laidbackware/cf-healthy-plugin/internal/util/http"
)

// Channels must be initialized but empty 
func GetLogs(
	cliConnection plugin.CliConnection, c logHttp.Client, sourceID string,
	quitC chan bool, shutdownC, sigkillC, httpErrorC chan int, errC chan error,
	debugMode bool,
) {

	skipSSL, err := cliConnection.IsSSLDisabled()
	log := log.New(os.Stderr, "", 0)
	if err != nil {
		log.Fatal(err)
		errC <- err
		return
	}

	hasAPI, err := cliConnection.HasAPIEndpoint()
	if err != nil {
		log.Fatalf("%s", err)
		errC <- err
		return
	}

	if !hasAPI {
		log.Fatalf("No API endpoint targeted.")
		errC <- err
		return
	}

	cfApiEndpoint, err := cliConnection.ApiEndpoint()
	if err != nil {
		log.Fatalf("%s", err)
		errC <- err
		return
	}

	logCacheAddr := strings.Replace(cfApiEndpoint, "api", "log-cache", 1)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: skipSSL,
	}
	c = logHttp.NewTokenClient(c, func() string {
		token, err := cliConnection.AccessToken()
		if err != nil {
			log.Fatalf("Unable to get Access Token: %s", err)
			errC <- err
		}
		return token
	})

	client := logcache.NewClient(logCacheAddr, logcache.WithHTTPClient(c))

	logcache.Walk(
		context.Background(),
		sourceID,
		logcache.Visitor(func(envelopes []*loggregator_v2.Envelope) bool {
			for _, e := range envelopes {
				select {
				case <-quitC:
					return false
				default:
					processMessage(shutdownC, sigkillC, httpErrorC, string(e.GetLog().GetPayload()), debugMode, log)
				}
			}
			return true
		}),
		client.Read,
		logcache.WithWalkEnvelopeTypes(logcache_v1.EnvelopeType_LOG),
		logcache.WithWalkStartTime(time.Unix(0, time.Now().Add(-5 * time.Second).UnixNano())),
		logcache.WithWalkBackoff(logcache.NewAlwaysRetryBackoff(250*time.Millisecond)),
	)
} 

func processMessage(shutdownC, sigkillC, httpErrorC chan int, logMessage string, debugMode bool, log Logger) {
	switch {
	case strings.Contains(logMessage, "successfully destroyed container for instance"):
		shutdownC <- getIntNonBlock(shutdownC) + 1
		if debugMode {
			log.Printf(logMessage)
		}
	case strings.Contains(logMessage, "Exit status 137 (exceeded 10s graceful shutdown interval)"):
		sigkillC <- getIntNonBlock(sigkillC) + 1
		log.Printf(logMessage)
	case strings.Contains(logMessage, "endpoint_failure"):
		httpErrorC <- getIntNonBlock(httpErrorC) + 1
		log.Printf("HTTP endpint error due to requests outstanding!")
	case (strings.Contains(logMessage, "stopping") || strings.Contains(logMessage, "destroying") ||
		strings.Contains(logMessage, "successfully") || strings.Contains(logMessage, "creating")):
		if debugMode {
			log.Printf(logMessage)
		}
	default:
		// enable debug logging for related events
		if debugMode {
			log.Printf(logMessage)
		}
	}
}
