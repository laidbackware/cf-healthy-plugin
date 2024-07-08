package collect_data

import (
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"context"
)

func AppLookup (cf *client.Client) (map[string]string, map[string]string, error) {
	appLookup := make(map[string]string)
	appSpaceLookup := make(map[string]string)
	apps, err := cf.Applications.ListAll(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range apps {
		appLookup[app.GUID] = app.Name
		appSpaceLookup[app.GUID] = app.Relationships.Space.Data.GUID
	}
	return appLookup, appSpaceLookup, nil
}

func SpaceLookup (cf *client.Client) (map[string]string, map[string]string, error) {
	spaceLookup := make(map[string]string)
	spaceOrgLookup := make(map[string]string)
	spaces, err := cf.Spaces.ListAll(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}
	for _, space := range spaces {
		spaceLookup[space.GUID] = space.Name
		spaceOrgLookup[space.GUID] = space.Relationships.Organization.Data.GUID
	}
	return spaceLookup, spaceOrgLookup, err
}

func OrgNameLookup (cf *client.Client) (map[string]string, error) {
	orgLookup := make(map[string]string)
	orgs, err := cf.Organizations.ListAll(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	for _, org := range orgs {
		orgLookup[org.GUID] = org.Name
	}
	return orgLookup, nil
}
