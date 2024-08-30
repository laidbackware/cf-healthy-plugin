package sig_check

import (
	"crypto/tls"
	"context"
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
	shutdown := make(chan int, 1)
	sigkill := make(chan int, 1)
	errChan := make(chan error, 1)

	go GetLogs(cli, http.DefaultClient, appGUID, quit, shutdown, sigkill, errChan, debugMode)
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
	
	shutdownInstances := getIntNonBlock(shutdown)
	sigkillInstances := getIntNonBlock(sigkill)
	log.Printf("Restart complete for app with guid: %s", appGUID)
	log.Printf("%d instaces restarted", shutdownInstances)
	if sigkillInstances > 0 {
		return fmt.Errorf("%d apps terminated using SIGKILL", sigkillInstances)
	}

	// non-blocking send quit
	select {
	case quit <- true:
	default:
	}
	
	return err
}

func restartApp(cf *client.Client, ctx context.Context, appGUID string, log Logger) (error) {
	dropletGUID, err := getCurrentDroplet(cf, ctx, appGUID)
	if err != nil {
		return err
	}
	log.Printf("Rolling restarting app with guid: %s", appGUID)
	log.Printf("Once each new instances passes it's health check an existing instance will be sent a SIGTERM signal")
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
	
	log.Printf("New instances running for %s. Waiting 30 seconds for shutdowns to complete", appGUID)
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
