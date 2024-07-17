package render_output

import (
	"testing"

	// "github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/laidbackware/cf-healthy-plugin/internal/collect_data"
	"github.com/stretchr/testify/assert"
)

func TestBuildTableArray(t *testing.T) {
	sheetContents := map[string]map[string]map[string][]collect_data.Process{
		"o1": {
			"s1": {
				"a1": {
					collect_data.Process{
						Type: "web",
						Instances: 1,
						AppGuid: "1-2",
						// HealthCheck: &resource.ProcessHealthCheck{
						// 	Type: "web",
						// 	Data: &resource
						// 	resource.ToOneRelationship{
						// 		Data: &resource.Relationship{
						// 			GUID: "a1uid",
						// 		},
						// 	},
						// },
					},
					collect_data.Process{
						Type: "worker",
						Instances: 1,
						AppGuid: "1-2",
					},
				},
			},
		},
		"o2": {
			"s2": {
				"a2": {
					collect_data.Process{
						Type: "web",
						Instances: 1,
						AppGuid: "1-2",
					},
				},
			},
			"s3": {
				"a3": {
					collect_data.Process{
						Type: "web",
						Instances: 1,
						AppGuid: "1-2",
					},
				},
				"a4": {
					collect_data.Process{
						Type: "web",
						Instances: 1,
						AppGuid: "1-2",
					},
				},
			},
		},
	}
	tableArray := buildTableArray(sheetContents)
	assert.Equal(t, len(tableArray), 5)
}