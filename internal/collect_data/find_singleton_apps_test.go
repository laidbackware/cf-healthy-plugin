package collect_data

import (
	// "fmt"
	// "os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

func TestFindSingletonApps(t *testing.T) {
	cf := initClient(t)
	healthState, err := CollectHealthState(cf)
	assert.Nil(t, err)
	assert.Greater(t, len(healthState.SingletonApps), 1)
}

func TestIterateProcesses(t *testing.T) {
	lookups := LookupTables{
		appNameLookup: map[string]string{"a1uid":"a1name", "a2uid":"a2name"},
		appSpaceLookup: map[string]string{"a1uid":"s1uid", "a2uid":"s2uid"},
		spaceNameLookup: map[string]string{"s1uid":"s1name", "s2uid":"s2name"},
		spaceOrgLookup: map[string]string{"s1uid":"o1uid", "s2uid":"o2uid"},
		orgNameLookup: map[string]string{"o1uid":"o1name", "o2uid":"o2name"},
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
		},
	}

	healthState, err := iterateProcesses(lookups, processes)
	assert.Nil(t, err)
	assert.Equal(t, len(healthState.SingletonApps), 2)
	assert.Equal(t, len(healthState.SingletonApps["o1name"]), 1)
	assert.Equal(t, len(healthState.SingletonApps["o1name"]["s1name"]), 1)
	assert.Equal(t, len(healthState.SingletonApps["o1name"]["s1name"]["a1name"]), 1)
}