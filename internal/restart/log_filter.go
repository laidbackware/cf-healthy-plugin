package restart

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	logHttp "github.com/laidbackware/cf-healthy-plugin/internal/util/http"
	logcache "code.cloudfoundry.org/go-log-cache/v3"
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

	// user, err := cliConnection.Username()
	// if err != nil {
	// 	log.Fatalf("%s", err)
	// }

	// org, err := cli.GetCurrentOrg()
	// if err != nil {
	// 	log.Fatalf("%s", err)
	// }

	// space, err := cli.GetCurrentSpace()
	// if err != nil {
	// 	log.Fatalf("%s", err)
	// }

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

	// walkStartTime := time.Now().Add(-5 * time.Second).UnixNano()


	envelopes, err := client.Read(
		context.Background(),
		sourceID,
		time.Now().Add(-5 * time.Second),
		// logcache.WithEndTime(o.endTime),
		// logcache.WithEnvelopeTypes(),
		// logcache.WithLimit(100),
		// logcache.WithDescending(),
		// logcache.WithNameFilter(o.nameFilter),
	)

	log.Print(len(envelopes))

	return err
}