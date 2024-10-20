package command

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"log"
	"os"
)

// HealthyPlugin is the struct implementing the interface defined by the core CLI. It can
// be found at "code.cloudfoundry.org/cli/plugin/plugin.go"
type HealthyPlugin struct{}

// Run must be implemented by any plugin because it is part of the
// plugin interface defined by the core CLI.
//
// Run(....) is the entry point when the core CLI is invoking a command defined
// by a plugin. The first parameter, plugin.CliConnection, is a struct that can
// be used to invoke cli commands. The second parameter, args, is a slice of
// strings. args[0] will be the name of the command, and will be followed by
// any additional arguments a cli user typed in.
//
// Any error handling should be handled with the plugin itself (this means printing
// user facing errors). The CLI will exit 0 if the plugin exits 0 and will exit
// 1 should the plugin exits nonzero.
func (c *HealthyPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	l := log.New(os.Stderr, "", 0)

	switch args[0] {
	case "health-report":
		generateHealthReport(cliConnection, args[1:], l)
	case "sig-check":
		sigCheck(cliConnection, args[1:], l)
	case "CLI-MESSAGE-UNINSTALL":
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unsupported command %s\n", args[0])
		os.Exit(1)
	}
}

// GetMetadata must be implemented as part of the plugin interface
// defined by the core CLI.
//
// GetMetadata() returns a PluginMetadata struct. The first field, Name,
// determines the name of the plugin which should generally be without spaces.
// If there are spaces in the name a user will need to properly quote the name
// during uninstall otherwise the name will be treated as separate arguments.
// The second value is a slice of Command structs. Our slice only contains one
// Command Struct, but could contain any number of them. The first field Name
// defines the command `cf basic-plugin-command` once installed into the CLI. The
// second field, HelpText, is used by the core CLI to display help information
// to the user in the core commands `cf help`, `cf`, or `cf -h`.
func (c *HealthyPlugin) GetMetadata() plugin.PluginMetadata {
	reportOptions := map[string]string{
		"--output, -o": "The output file, with or without path.",
		"--format, -f": "The format of the output file. (json, xlsx).",
	}
	sigOptions := map[string]string{
		"--debug, -d": "Enabled debug logging, which will desplay more logs.",
	}

	return plugin.PluginMetadata{
		Name: "HealthyPlugin",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 3,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 7,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "health-report",
				HelpText: "Generate report on app to find singletons and port health checks",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage:   "cf health-report [OPTIONS]",
					Options: reportOptions,
				},
			},
			{
				Name:     "sig-check",
				HelpText: "Rolling restart app and check that no SIGKILLs were sent",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage:   "cf sig-check [OPTIONS] <app-name>",
					Options: sigOptions,
				},
			},
		},
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
