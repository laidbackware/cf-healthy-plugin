package restart

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache/v3"
	logcache_v1 "code.cloudfoundry.org/go-log-cache/v3/rpc/logcache_v1"
	"code.cloudfoundry.org/go-loggregator/v10/rpc/loggregator_v2"
	logHttp "github.com/laidbackware/cf-healthy-plugin/internal/util/http"

	// "github.com/cloudfoundry/go-cfclient/v3/client"
	// "github.com/cloudfoundry/go-cfclient/v3/config"
	"code.cloudfoundry.org/cli/plugin"
	// "github.com/cloudfoundry/go-cfclient/v3/resource"
)

func GetLogs(cliConnection plugin.CliConnection, c logHttp.Client, sourceID string) error {
	skipSSL, err := cliConnection.IsSSLDisabled()
	log := log.New(os.Stderr, "", 0)
	if err != nil {
		log.Fatal(err)
	}
	
	hasAPI, err := cliConnection.HasAPIEndpoint()
	if err != nil {
		log.Fatalf("%s", err)
	}

	if !hasAPI {
		log.Fatalf("No API endpoint targeted.")
	}

	tokenURL, err := cliConnection.ApiEndpoint()
	if err != nil {
		log.Fatalf("%s", err)
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

	var startTime int64 = 1723125134
	
	envelopes, err := client.Read(
		context.Background(),
		sourceID,
		time.Unix(startTime, 0),
		// time.Now().Add(-120 * time.Second),
		// logcache.WithEndTime(o.endTime),
		logcache.WithEnvelopeTypes(logcache_v1.EnvelopeType_LOG),
		// logcache.WithLimit(100),
		// logcache.WithDescending(),
		// logcache.WithNameFilter(o.nameFilter),
	)
	log.Print(len(envelopes))


	var shutdown 	int
	var sigkill 	int

	logcache.Walk(
		context.Background(),
		sourceID,
		logcache.Visitor(func(envelopes []*loggregator_v2.Envelope) bool {
			for _, e := range envelopes {
				processMessage(&shutdown, &sigkill, string(e.GetLog().GetPayload()))
			}
			return true
		}),
		client.Read,
		// logcache.WithWalkStartTime(time.Unix(0, walkStartTime)),
		logcache.WithWalkEnvelopeTypes(logcache_v1.EnvelopeType_LOG),
		logcache.WithWalkBackoff(logcache.NewAlwaysRetryBackoff(250*time.Millisecond)),
		// logcache.WithWalkNameFilter(o.nameFilter),
	)
	return err
}

func processMessage(shutdown, sigkill *int, logMessage string) {
	switch {
	case strings.Contains(logMessage, "successfully destroyed container for instance"):
		*shutdown ++
	case strings.Contains(logMessage, "Exit status 137 (exceeded 10s graceful shutdown interval)"):
		*sigkill++
	}
	fmt.Println(logMessage)
}