package render_output

import (
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/resource"
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
						HealthCheck: &resource.ProcessHealthCheck{
							Type: "port",
						},
					},
					collect_data.Process{
						Type: "worker",
						Instances: 1,
						AppGuid: "1-2",
						HealthCheck: &resource.ProcessHealthCheck{
							Type: "http",
							Data: resource.ProcessData{
								InvocationTimeout: createIntPointer(30),
								Endpoint: createStringPointer("blah"),
							},
						},
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
						HealthCheck: &resource.ProcessHealthCheck{
							Type: "http",
							Data: resource.ProcessData{
								InvocationTimeout: createIntPointer(10),
								Endpoint: createStringPointer("blah"),
							},
						},
					},
				},
			},
			"s3": {
				"a3": {
					collect_data.Process{
						Type: "web",
						Instances: 1,
						AppGuid: "1-2",
						HealthCheck: &resource.ProcessHealthCheck{
							Type: "port",
						},
					},
				},
				"a4": {
					collect_data.Process{
						Type: "web",
						Instances: 1,
						AppGuid: "1-2",
						HealthCheck: &resource.ProcessHealthCheck{
							Type: "port",
						},
					},
				},
			},
		},
	}
	tableArray := buildTableArray(sheetContents, true)
	assert.Equal(t, len(tableArray), 5)
}

func createIntPointer(x int) *int {
	return &x
}

func createStringPointer(s string) *string {
	return &s
}