package restart

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

func RestartApp(cf *client.Client, appGuid string, dropletGuid string) error {
	// TODO check if there is currently an active deployment

	c := resource.NewDeploymentCreate(appGuid)
	c.Droplet = &resource.Relationship{
		GUID: dropletGuid,
	}

	dep, err := cf.Deployments.Create(context.Background(), c)
	if err != nil {
		return err
	}

	for {
		dep, err = cf.Deployments.Get(context.Background(), dep.GUID)
		// app, err := cf.Applications.Get(context.Background(), appGuid)
		if err != nil {
			return err
		}
		if dep.Status.Value == "FINALIZED" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Println("")

	return err
}
