package collect_data

import (
	// "fmt"
	// "os"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/stretchr/testify/assert"
)

func TestFindSingletonApps(t *testing.T) {
	cf := initClient(t)
	healthState, err := CollectHealthState(cf)
	assert.Nil(t, err)
	assert.Greater(t, len(healthState.SingletonApps), 0)
}

func TestIterateProcesses(t *testing.T) {
	lookups := LookupTables{
		appNameLookup:      map[string]string{"a1uid": "a1name", "a2uid": "a2name", "a3uid": "a3name", "a4uid": "a4name"},
		appSpaceNameLookup: map[string]string{"a1uid": "s1name", "a2uid": "s2name", "a3uid": "s1name", "a4uid": "s2name"},
		appOrgNameLookup:   map[string]string{"a1uid": "o1name", "a2uid": "o2name", "a3uid": "o1name", "a4uid": "o2name", "system-app": "system"},
	}

	processes := []*resource.Process{
		{
			Relationships: resource.ProcessRelationships{
				App: resource.ToOneRelationship{
					Data: &resource.Relationship{
						GUID: "a1uid",
					},
				},
			},
			Instances: 1,
			HealthCheck: resource.ProcessHealthCheck{
				Type: "port",
			},
		},
		{
			Relationships: resource.ProcessRelationships{
				App: resource.ToOneRelationship{
					Data: &resource.Relationship{
						GUID: "a2uid",
					},
				},
			},
			Instances: 1,
			HealthCheck: resource.ProcessHealthCheck{
				Type: "http",
				Data: resource.ProcessHealthCheckData{
					Interval: createIntPointer(30),
				},
			},
		},
		{
			Relationships: resource.ProcessRelationships{
				App: resource.ToOneRelationship{
					Data: &resource.Relationship{
						GUID: "a3uid",
					},
				},
			},
			Instances: 2,
			HealthCheck: resource.ProcessHealthCheck{
				Type: "http",
				Data: resource.ProcessHealthCheckData{
					Interval: createIntPointer(10),
				},
			},
		},
	}

	healthState, err := iterateProcesses(lookups, processes)
	assert.Nil(t, err)
	assert.Equal(t, len(healthState.SingletonApps), 2)
	assert.Equal(t, len(healthState.SingletonApps["o1name"]["s1name"]["a1name"]), 1)
	assert.Equal(t, len(healthState.PortHealthCheck), 1)
	assert.Equal(t, len(healthState.PortHealthCheck["o1name"]["s1name"]["a1name"]), 1)
	assert.Equal(t, len(healthState.LongInterval), 2)
	assert.Equal(t, len(healthState.LongInterval["o2name"]["s2name"]["a2name"]), 1)
}

func createIntPointer(x int) *int {
	return &x
}
