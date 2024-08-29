package sig_check

import (
	"context"
	"time"

	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

func RestartApp(cf *client.Client, appGUID string, log Logger) error {
	// TODO check if there is currently an active deployment

	ctx := context.Background()

	dropletGUID, err := getCurrentDroplet(cf, ctx, appGUID)
	if err != nil {
		return err
	}

	log.Printf("Restarting app with guid: %s", appGUID)

	c := resource.NewDeploymentCreate(appGUID)
	c.Droplet = &resource.Relationship{
		GUID: dropletGUID,
	}

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

	
	log.Printf("New instances running for %s. Wating 15 seconds for shutdowns to complete", appGUID)
	time.Sleep(15 * time.Second)

	log.Printf("Restart complete for app with guid: %s", appGUID)

	return err
}

func getCurrentDroplet(cf *client.Client, ctx context.Context, appGUID string) (string, error) {
	droplet, err := cf.Droplets.GetCurrentAssociationForApp(ctx, appGUID)
	if err != nil {
		return "", err
	}
	return droplet.Data.GUID, nil
}
