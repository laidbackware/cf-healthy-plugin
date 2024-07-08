package collect_data

import (
	"context"
	// "fmt"

	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

type LookupTables struct {
	appNameLookup			map[string]string
	appSpaceLookup  	map[string]string
	spaceNameLookup 	map[string]string
	spaceOrgLookup		map[string]string
	orgNameLookup			map[string]string
}

func FindSingletonApps(cf *client.Client) (map[string]map[string]map[string][]*resource.Process, error) {
	var lookups LookupTables
	
	var err error
	lookups.appNameLookup, lookups.appSpaceLookup, err = AppLookup(cf)
	if err != nil {
		return nil, err
	}
	lookups.spaceNameLookup, lookups.spaceOrgLookup, err = SpaceLookup(cf)
	if err != nil {
		return nil, err
	}
	lookups.orgNameLookup, err = OrgNameLookup(cf)
	if err != nil {
		return nil, err
	}
	processes, err := cf.Processes.ListAll(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	singletonApps, err := iterateProcesses(lookups, processes)

	return singletonApps, err
}


func iterateProcesses(lookups LookupTables, processes []*resource.Process) (map[string]map[string]map[string][]*resource.Process, error) {
	var err error
	singletonApps := make(map[string]map[string]map[string][]*resource.Process)

	for _, process := range processes{
		if process.Instances < 2 && process.Type != "task" {
			appName := lookups.appNameLookup[process.Relationships.App.Data.GUID]
			appSpace := lookups.spaceNameLookup[
				lookups.appSpaceLookup[process.Relationships.App.Data.GUID]]
			appOrg := lookups.orgNameLookup[
				lookups.spaceOrgLookup[lookups.appSpaceLookup[process.Relationships.App.Data.GUID]]]

			if _, ok := singletonApps[appOrg]; !ok {
				singletonApps[appOrg] = make(map[string]map[string][]*resource.Process)
			}
			if _, ok := singletonApps[appOrg][appSpace]; !ok {
				singletonApps[appOrg][appSpace] = make(map[string][]*resource.Process)
			}
			if _, ok := singletonApps[appOrg][appSpace][appName]; !ok {
				singletonApps[appOrg][appSpace][appName] = []*resource.Process{process}
			} else  {
				singletonApps[appOrg][appSpace][appName] = append(singletonApps[appOrg][appSpace][appName], process)
			}
		}
	}
	return singletonApps, err
}

