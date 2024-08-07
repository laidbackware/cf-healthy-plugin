package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/cli/cf/flags"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/laidbackware/cf-healthy-plugin/internal/collect_data"
	"github.com/laidbackware/cf-healthy-plugin/internal/render_output"
)

func generateHealthReport(cliConnection plugin.CliConnection, args []string) {
	fc, err := parseArguements(args)
	handleError(err)

	outputFile := strings.ToLower(fc.String("output"))
	fileFormat := strings.ToLower(fc.String("format"))

	if fileFormat != "json" && fileFormat != "xlsx" {
		fmt.Fprintln(os.Stderr, "Requested output format is invlaid. Please use: [json, xlsx]")
		os.Exit(1)
	}

	if outputFile == "" {
		currentDir, err := os.Getwd()
		handleError(err)
		if fileFormat == "json" {
			outputFile = filepath.Join(currentDir, "report.json")
		} else {
			outputFile = filepath.Join(currentDir, "report.xlsx")
		}
	}

	cf, err := createCFClient(cliConnection)
	handleError(err)

	// var healthState collect_data.HealthState
	healthState, err := collect_data.CollectHealthState(cf)
	handleError(err)

	switch fileFormat {
	case "xlsx":
		handleError(render_output.WriteSheet(healthState, outputFile))
	case "json":
		handleError(render_output.WriteJSON(healthState, outputFile))
	default:
		fmt.Fprintf(os.Stderr, "File format %s is not support. Please use [json, xlsx]\n", fileFormat)
		os.Exit(1)
	}
	fmt.Printf("Written file: %s\n", outputFile)
}

func parseArguements(args []string) (flags.FlagContext, error) {
	fc := flags.New()
	fc.NewStringFlag("output", "o", "The output file, with or without path.")
	fc.NewStringFlagWithDefault("format", "f", "The format of the output file. (json, xlsx).", "xlsx")
	err := fc.Parse(args...)
	return fc, err
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
