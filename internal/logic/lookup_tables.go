package logic

import (
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"context"
	"fmt"
	"os"
)

func SpaceNameLookup (cf *client.Client) map[string]string {
	spaceLookup := make(map[string]string)
	spaces, err := cf.Spaces.ListAll(context.Background(), nil)
	handleError(err)
	for _, space := range spaces {
		spaceLookup[space.GUID] = space.Name
	}
	return spaceLookup
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}