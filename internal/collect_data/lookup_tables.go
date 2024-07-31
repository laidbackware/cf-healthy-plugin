package collect_data

import (
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"context"
)

func AppLookup (cf *client.Client, spaceOrgNameLookup map[string]string,
		spaceLookup map[string]string) (map[string]string, map[string]string, 
		map[string]string, error) {
	appNameLookup := make(map[string]string)
	appSpaceNameLookup := make(map[string]string)
	appOrgNameLookup := make(map[string]string)
	apps, err := cf.Applications.ListAll(context.Background(), nil)
	if err != nil {
		return nil, nil, nil, err
	}
	for _, app := range apps {
		appNameLookup[app.GUID] = app.Name
		appSpaceNameLookup[app.GUID] = spaceLookup[app.Relationships.Space.Data.GUID]
		appOrgNameLookup[app.GUID] = spaceOrgNameLookup[app.Relationships.Space.Data.GUID]
	}
	return appNameLookup, appSpaceNameLookup, appOrgNameLookup, nil
}

func SpaceLookup (cf *client.Client, orgLookup map[string]string) (
		map[string]string, map[string]string, error) {
	spaceNameLookup := make(map[string]string)
	spaceOrgNameLookup := make(map[string]string)
	spaces, err := cf.Spaces.ListAll(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}
	for _, space := range spaces {
		spaceNameLookup[space.GUID] = space.Name
		spaceOrgNameLookup[space.GUID] = orgLookup[space.Relationships.Organization.Data.GUID]
	}
	return spaceNameLookup, spaceOrgNameLookup, nil
}

func OrgNameLookup (cf *client.Client) (map[string]string, error) {
	orgNameLookup := make(map[string]string)
	orgs, err := cf.Organizations.ListAll(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	for _, org := range orgs {
		orgNameLookup[org.GUID] = org.Name
	}
	return orgNameLookup, nil
}
