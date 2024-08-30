package sig_check

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache/v3"
	logcache_v1 "code.cloudfoundry.org/go-log-cache/v3/rpc/logcache_v1"
	"code.cloudfoundry.org/go-loggregator/v10/rpc/loggregator_v2"
	logHttp "github.com/laidbackware/cf-healthy-plugin/internal/util/http"
	"code.cloudfoundry.org/cli/plugin"
)

func GetLogs(
		cliConnection plugin.CliConnection, c logHttp.Client, sourceID string, 
		quit chan bool, shutdown, sigkill chan int, errC chan error,
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

	tokenURL, err := cliConnection.ApiEndpoint()
	if err != nil {
		log.Fatalf("%s", err)
		errC <- err
		return
	}

	logCacheAddr := strings.Replace(tokenURL, "api", "log-cache", 1)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: skipSSL,
	}
	c = logHttp.NewTokenClient(c, func() string {
		token, err := cliConnection.AccessToken()
		if err != nil {
			log.Fatalf("Unable to get Access Token: %s", err)
		}
		return token
	})

	client := logcache.NewClient(logCacheAddr, logcache.WithHTTPClient(c))

	walkStartTime := time.Now().Add(-5 * time.Second).UnixNano()

	logcache.Walk(
		context.Background(),
		sourceID,
		logcache.Visitor(func(envelopes []*loggregator_v2.Envelope) bool {
			for _, e := range envelopes {
				select {
        case <- quit:
					return false
        default:
					processMessage(shutdown, sigkill, string(e.GetLog().GetPayload()), debugMode, log)
        }
			}
			return true
		}),
		client.Read,
		logcache.WithWalkEnvelopeTypes(logcache_v1.EnvelopeType_LOG),
		logcache.WithWalkStartTime(time.Unix(0, walkStartTime)),
		logcache.WithWalkBackoff(logcache.NewAlwaysRetryBackoff(250*time.Millisecond)),
	)
}

func processMessage(shutdown, sigkill chan int, logMessage string, debugMode bool, log Logger) {
	switch {
	case strings.Contains(logMessage, "successfully destroyed container for instance"):
		// current := 
		// current := <- shutdown
		// current++
		debugLog(logMessage, debugMode, log)
		sendIntNonBlock(shutdown, getIntNonBlock(shutdown) + 1)
		// shutdown <- currente
	case strings.Contains(logMessage, "Exit status 137 (exceeded 10s graceful shutdown interval)"):
		// current := <- sigkill
		// current++
		// sigkill <- current

		debugLog(logMessage, debugMode, log)
		sendIntNonBlock(sigkill, getIntNonBlock(sigkill) + 1)
	// enable debug logging for related events
	case (
		strings.Contains(logMessage, "stopping") || strings.Contains(logMessage, "destroying") || 
		strings.Contains(logMessage, "successfully") || strings.Contains(logMessage, "creating")): 
		debugLog(logMessage, debugMode, log)
	// default: 
	// 	debugLog(logMessage, debugMode)
	}
}

func debugLog(logMessage string, debugMode bool, log Logger) {
	if debugMode{
		log.Printf(logMessage)
	}
}