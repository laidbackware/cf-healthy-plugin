package render_output

import (
	"encoding/json"
	"os"

	"github.com/laidbackware/cf-healthy-plugin/internal/collect_data"
)

func WriteJSON(healthState collect_data.HealthState, outputFile string) (err error) {
	
	outputJson, err := json.MarshalIndent(healthState, "", "    ")
	if err != nil {
		return
	}
  err = os.WriteFile(outputFile, outputJson, 0644)
	return
}