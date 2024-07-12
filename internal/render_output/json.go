package render_output

import (
	"encoding/json"
	"os"

	"github.com/cloudfoundry/go-cfclient/v3/resource"

)

func WriteJSON(singletonApps map[string]map[string]map[string][]*resource.Process, outputFile string) (err error) {
	
	outputJson, err := json.MarshalIndent(singletonApps, "", "    ")
	if err != nil {
		return
	}
  err = os.WriteFile(outputFile, outputJson, 0644)
	return
}