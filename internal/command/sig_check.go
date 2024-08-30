package command

import (
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	flags "github.com/jessevdk/go-flags"
	"github.com/laidbackware/cf-healthy-plugin/internal/sig_check"
)

func sigCheck(cli plugin.CliConnection, args []string, log Logger) {
	cf, err := createCFClient(cli)
	handleError(err)

	o, err := newSigOptions(cli, args)
	handleError(err)

	err = sig_check.SigCheck(cli, cf, o.appGUID, log, o.debugMode)
	handleError(err)
}

type sigOptions struct {
	appGUID   string
	debugMode bool
}

type sigOptionFlags struct {
	Debug   bool  `bool:"debug" short:"d"`
}

func newSigOptions(cli plugin.CliConnection, args []string) (sigOptions, error) {
	opts := sigOptionFlags{}

	args, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return sigOptions{}, err
	}

	if len(args) != 1 {
		return sigOptions{}, fmt.Errorf("expected app name as argument, got %d", len(args))
	}

	appGUID, err := getAppGUID(args[0], cli)
	if err != nil {
		return sigOptions{}, err
	}

	o := sigOptions{
		appGUID:   appGUID,
		debugMode: opts.Debug,
	}

	return o, nil
}

func getAppGUID(appName string, cli plugin.CliConnection) (string, error) {
	r, err := cli.CliCommandWithoutTerminalOutput(
		"app",
		appName,
		"--guid",
	)

	if err == nil {
		return strings.Join(r, ""), nil
	}
	return "", err
}
