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

// All return data stored as a nested map to allow for it to be easily sorted before rendering
type HealthState struct {
	SingletonApps 		map[string]map[string]map[string][]Process `json:"singleton_apps"`
	PortHealthCheck 	map[string]map[string]map[string][]Process `json:"port_health_checks"`
	DefaultHttpTime		map[string]map[string]map[string][]Process `json:"default_http_interval"`
	AllProcesses			map[string]map[string]map[string][]Process `json:"all_process_data"`
}

type Process struct {
	Instances 				int														`json:"instance_count"`
	Type 							string												`json:"process_type"`
	AppGuid						string												`json:"app_guid"`
	HealthCheck				*resource.ProcessHealthCheck	`json:"health_check"`
}

func CollectHealthState(cf *client.Client) (HealthState, error) {
	var lookups LookupTables
	
	var err error
	lookups.appNameLookup, lookups.appSpaceLookup, err = AppLookup(cf)
	if err != nil {
		return HealthState{}, err
	}
	lookups.spaceNameLookup, lookups.spaceOrgLookup, err = SpaceLookup(cf)
	if err != nil {
		return HealthState{}, err
	}
	lookups.orgNameLookup, err = OrgNameLookup(cf)
	if err != nil {
		return HealthState{}, err
	}
	processes, err := cf.Processes.ListAll(context.Background(), nil)
	if err != nil {
		return HealthState{}, err
	}

	healthState, err := iterateProcesses(lookups, processes)

	return healthState, err
}


func iterateProcesses(lookups LookupTables, processes []*resource.Process) (HealthState, error) {
	var err error
	healthState := HealthState {
		SingletonApps: 		make(map[string]map[string]map[string][]Process),
		PortHealthCheck: 	make(map[string]map[string]map[string][]Process),
		DefaultHttpTime: 	make(map[string]map[string]map[string][]Process),
		AllProcesses: 		make(map[string]map[string]map[string][]Process),
	}

	for _, fullProcess := range processes{
		process := Process {
			AppGuid: 			fullProcess.Relationships.App.Data.GUID,
			Instances: 		fullProcess.Instances,
			Type: 				fullProcess.Type,
			HealthCheck:  &fullProcess.HealthCheck,
		}
		addProcess(lookups, healthState.AllProcesses, process)
		if process.Instances < 2 && process.Type != "task" {
				addProcess(lookups, healthState.SingletonApps, process)
		}
		if process.HealthCheck.Type == "port" {
			addProcess(lookups, healthState.PortHealthCheck, process)
		}
		if process.HealthCheck.Data.InvocationTimeout != nil && *process.HealthCheck.Data.InvocationTimeout == 30 {
			addProcess(lookups, healthState.DefaultHttpTime, process)
		}	
	}
	return healthState, err
}

func addProcess(lookups LookupTables, targetMap map[string]map[string]map[string][]Process, process Process) {
	appName := lookups.appNameLookup[process.AppGuid]
	appSpace := lookups.spaceNameLookup[
		lookups.appSpaceLookup[process.AppGuid]]
	appOrg := lookups.orgNameLookup[
		lookups.spaceOrgLookup[lookups.appSpaceLookup[process.AppGuid]]]

	if _, ok := targetMap[appOrg]; !ok {
		targetMap[appOrg] = make(map[string]map[string][]Process)
	}
	if _, ok := targetMap[appOrg][appSpace]; !ok {
		targetMap[appOrg][appSpace] = make(map[string][]Process)
	}
	if _, ok := targetMap[appOrg][appSpace][appName]; !ok {
		targetMap[appOrg][appSpace][appName] = []Process{process}
	} else  {
		targetMap[appOrg][appSpace][appName] = append(targetMap[appOrg][appSpace][appName], process)
	}
}