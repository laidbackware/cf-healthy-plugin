package healthy_plugin

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/laidbackware/cf-healthy-plugin/internal/collect_data"
	"github.com/laidbackware/cf-healthy-plugin/internal/sheet_writer"
)

func RunPlugin(cliConnection plugin.CliConnection) {
	cf, err := createCFClient(cliConnection)
	handleError(err)

	singletonApps, err := collect_data.FindSingletonApps(cf)
	handleError(err)

	outputFile, err := sheet_writer.WriteSheet(singletonApps, "tests.xlsx")
	handleError(err)
	fmt.Printf("Written file: %s\n", outputFile)
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
