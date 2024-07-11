package command

import (
	"fmt"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/cf/flags"
	"github.com/laidbackware/cf-healthy-plugin/internal/collect_data"
	"github.com/laidbackware/cf-healthy-plugin/internal/sheet_writer"
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

	handleError(sheet_writer.WriteSheet(singletonApps, outputFile))
	fmt.Printf("Written file: %s\n", outputFile)
}

func parseArguements(args []string) (flags.FlagContext, error) {
	fc := flags.New()
	fc.NewStringFlag("output", "o", "The output file, with or without path.")
	fc.NewStringFlag("format", "f", "The format of the output file. (json, xlsx).")
	err := fc.Parse(args...)
	return fc, err
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
