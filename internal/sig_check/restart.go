package sig_check

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

func SigCheck(cli plugin.CliConnection, cf *client.Client, appGUID string, log Logger, debugMode bool) error {
	// TODO check if there is currently an active deployment

	ctx := context.Background()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true, //nolint:gosec
	}
	quit := make(chan bool, 1)
	shutdownC := make(chan int, 1)
	sigkillC := make(chan int, 1)
	httpErrorC := make(chan int, 1)
	errChan := make(chan error, 1)

	go GetLogs(cli, http.DefaultClient, appGUID, quit, shutdownC, sigkillC, httpErrorC, errChan, debugMode)
	errReceive := getErrNonBlock(errChan)
	if errReceive != nil {
		return errReceive
	}

	err := restartApp(cf, ctx, appGUID, log)
	if err != nil {
		return err
	}

	errReceive = getErrNonBlock(errChan)
	if errReceive != nil {
		return errReceive
	}

	shutdownCount := getIntNonBlock(shutdownC)
	sigkillCount := getIntNonBlock(sigkillC)
	endpointFailrureCount := getIntNonBlock(httpErrorC)
	log.Printf("Restart complete for app guid: %s", appGUID)
	log.Printf("%d instaces restarted", shutdownCount)

	// non-blocking send quit
	select {
	case quit <- true:
	default:
	}

	if sigkillCount > 0 || endpointFailrureCount > 0 {
		return fmt.Errorf("\n!!!FAILED!!!\n%d instances terminated using SIGKILL.\n%d http endpoint error/s encountered", 
			sigkillCount, endpointFailrureCount)
	}

	log.Printf("Success. All apps responded to SIGTERM.")
	return err
}

func restartApp(cf *client.Client, ctx context.Context, appGUID string, log Logger) error {
	dropletGUID, err := getCurrentDroplet(cf, ctx, appGUID)
	if err != nil {
		return err
	}
	log.Printf("Rolling restarting app guid: %s", appGUID)
	// TODO inform user of how many instances there will be

	c := resource.NewDeploymentCreate(appGUID)
	c.Droplet = &resource.Relationship{
		GUID: dropletGUID,
	}

	// Create deployment against existing droplet to rolling restart
	dep, err := cf.Deployments.Create(ctx, c)
	if err != nil {
		return err
	}

	// Run infinite loop to check deployment state
	for {
		dep, err = cf.Deployments.Get(ctx, dep.GUID)
		if err != nil {
			return err
		}
		if dep.Status.Value == "FINALIZED" {
			break
		}
		time.Sleep(2 * time.Second)
	}

	log.Printf("New instances running for: %s.\nWaiting 30 seconds for shutdowns to complete", appGUID)
	time.Sleep(30 * time.Second)

	return nil
}

func getCurrentDroplet(cf *client.Client, ctx context.Context, appGUID string) (string, error) {
	droplet, err := cf.Droplets.GetCurrentAssociationForApp(ctx, appGUID)
	if err != nil {
		return "", err
	}
	return droplet.Data.GUID, nil
}
