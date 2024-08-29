package command

import (
	"fmt"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/laidbackware/cf-healthy-plugin/internal/sig_check"
)

func sigCheck(cli plugin.CliConnection, args []string, log Logger){

	cf, err := createCFClient(cli)
	handleError(err)

	o, err := newSigOptions(cli, args, log)
	handleError(err)

	log.Printf(o.appGUID)

	sig_check.RestartApp(cf, o.appGUID, log)

}

type sigOptions struct {
	timeout     	int16
	appGUID 			string
}

type sigOptionFlags struct {
	Timeout     int16  `long:"timeout" short:"t"`
}

func newSigOptions(cli plugin.CliConnection, args []string, log Logger) (sigOptions, error) {
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
		timeout:	opts.Timeout,
		appGUID:	appGUID,
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