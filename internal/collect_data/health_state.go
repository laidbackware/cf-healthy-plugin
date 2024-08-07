package collect_data

import (
	"context"
	"slices"
	// "fmt"

	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

type LookupTables struct {
	appNameLookup      map[string]string
	appSpaceNameLookup map[string]string
	appOrgNameLookup   map[string]string
}

// All return data stored as a nested map to allow for it to be easily sorted before rendering
type HealthState struct {
	SingletonApps   map[string]map[string]map[string][]Process `json:"singleton_apps"`
	PortHealthCheck map[string]map[string]map[string][]Process `json:"port_health_checks"`
	LongInterval    map[string]map[string]map[string][]Process `json:"default_http_interval"`
	AllProcesses    map[string]map[string]map[string][]Process `json:"all_process_data"`
}

type Process struct {
	Instances   int                          `json:"instance_count"`
	Type        string                       `json:"process_type"`
	AppGuid     string                       `json:"app_guid"`
	HealthCheck *resource.ProcessHealthCheck `json:"health_check"`
}

func CollectHealthState(cf *client.Client) (HealthState, error) {
	var lookups LookupTables

	var err error
	orgNameLookup, err := OrgNameLookup(cf)
	if err != nil {
		return HealthState{}, err
	}
	var spaceNameLookup map[string]string
	var spaceOrgNameLookup map[string]string
	spaceNameLookup, spaceOrgNameLookup, err = SpaceLookup(cf, orgNameLookup)
	if err != nil {
		return HealthState{}, err
	}
	lookups.appNameLookup, lookups.appSpaceNameLookup, lookups.appOrgNameLookup, err = AppLookup(cf, spaceNameLookup, spaceOrgNameLookup)
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
	systemOrgs := []string{"system", "app-metrics-v2", "healthwatch2"}
	var err error
	healthState := HealthState{
		SingletonApps:   make(map[string]map[string]map[string][]Process),
		PortHealthCheck: make(map[string]map[string]map[string][]Process),
		LongInterval:    make(map[string]map[string]map[string][]Process),
		AllProcesses:    make(map[string]map[string]map[string][]Process),
	}

	for _, fullProcess := range processes {
		// Skip system orgs
		if slices.Contains(systemOrgs, lookups.appOrgNameLookup[fullProcess.Relationships.App.Data.GUID]) {
			continue
		}
		// Skip tasks
		if fullProcess.Type == "task" {
			continue
		}
		process := Process{
			AppGuid:     fullProcess.Relationships.App.Data.GUID,
			Instances:   fullProcess.Instances,
			Type:        fullProcess.Type,
			HealthCheck: &fullProcess.HealthCheck,
		}
		addProcess(lookups, healthState.AllProcesses, process)
		if process.Instances < 2 {
			addProcess(lookups, healthState.SingletonApps, process)
		}
		if process.HealthCheck.Type == "port" {
			addProcess(lookups, healthState.PortHealthCheck, process)
		}
		if process.HealthCheck.Data.Interval != nil && *process.HealthCheck.Data.Interval > 15 {
			addProcess(lookups, healthState.LongInterval, process)
		} else if process.HealthCheck.Data.Interval == nil {
			addProcess(lookups, healthState.LongInterval, process)
		}
	}
	return healthState, err
}

func addProcess(lookups LookupTables, targetMap map[string]map[string]map[string][]Process, process Process) {
	appName := lookups.appNameLookup[process.AppGuid]
	appSpace := lookups.appSpaceNameLookup[process.AppGuid]
	appOrg := lookups.appOrgNameLookup[process.AppGuid]

	if _, ok := targetMap[appOrg]; !ok {
		targetMap[appOrg] = make(map[string]map[string][]Process)
	}
	if _, ok := targetMap[appOrg][appSpace]; !ok {
		targetMap[appOrg][appSpace] = make(map[string][]Process)
	}
	if _, ok := targetMap[appOrg][appSpace][appName]; !ok {
		targetMap[appOrg][appSpace][appName] = []Process{process}
	} else {
		targetMap[appOrg][appSpace][appName] = append(targetMap[appOrg][appSpace][appName], process)
	}
}
