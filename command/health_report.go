package command

import (
	"fmt"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/cf/flags"
	"github.com/laidbackware/cf-healthy-plugin/internal/collect_data"
	"github.com/laidbackware/cf-healthy-plugin/internal/render_output"
)

func healthReport(cliConnection plugin.CliConnection, args []string) {
	fc, err := parseArguements(args)
	handleError(err)

	outputFile := fc.String("output")
  fileFormat := fc.String("format")
	fmt.Println(fileFormat)

	if outputFile == "" {
		currentDir, err := os.Getwd()
		handleError(err)
		outputFile = filepath.Join(currentDir, "report.xlsx")
	}
	
	cf, err := createCFClient(cliConnection)
	handleError(err)

	singletonApps, err := collect_data.FindSingletonApps(cf)
	handleError(err)

	switch fileFormat {
	case "xlsx":
		handleError(render_output.WriteSheet(singletonApps, outputFile))
	case "json":
		handleError(render_output.WriteJSON(singletonApps, outputFile))
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
