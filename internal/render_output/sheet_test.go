package render_output

import (
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/stretchr/testify/assert"
)

func TestBuildTableArray(t *testing.T) {
	singletonApps := map[string]map[string]map[string][]*resource.Process{
		"o1": {
			"s1": {
				"a1": {
					&resource.Process{
						Type: "web",
						Relationships: resource.ProcessRelationships{
							App: resource.ToOneRelationship{
								Data: &resource.Relationship{
									GUID: "a1uid",
								},
							},
						},
					},
					&resource.Process{
						Type: "worker",
						Relationships: resource.ProcessRelationships{
							App: resource.ToOneRelationship{
								Data: &resource.Relationship{
									GUID: "a1uid",
								},
							},
						},
					},
				},
			},
		},
		"o2": {
			"s2": {
				"a2": {
					&resource.Process{
						Type: "web",
						Relationships: resource.ProcessRelationships{
							App: resource.ToOneRelationship{
								Data: &resource.Relationship{
									GUID: "a2uid",
								},
							},
						},
					},
				},
			},
			"s3": {
				"a3": {
					&resource.Process{
						Type: "web",
						Relationships: resource.ProcessRelationships{
							App: resource.ToOneRelationship{
								Data: &resource.Relationship{
									GUID: "a3uid",
								},
							},
						},
					},
				},
				"a4": {
					&resource.Process{
						Type: "web",
						Relationships: resource.ProcessRelationships{
							App: resource.ToOneRelationship{
								Data: &resource.Relationship{
									GUID: "a4uid",
								},
							},
						},
					},
				},
			},
		},
	}
	tableArray := buildTableArray(singletonApps)
	assert.Equal(t, len(tableArray), 5)
}