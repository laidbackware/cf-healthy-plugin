package command

import (
	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/config"
	"strings"
)

func createCFClient(cliConnection plugin.CliConnection) (*client.Client, error) {
	u, err := cliConnection.ApiEndpoint()
	if err != nil {
		return nil, err
	}

	t, err := cliConnection.AccessToken()
	if err != nil {
		return nil, err
	}
	t = strings.TrimPrefix(t, "bearer ")
	
	skipSSLValidation, err := cliConnection.IsSSLDisabled()
	if err != nil {
		return nil, err
	}
	var cfg *config.Config
	if skipSSLValidation {
		cfg, err = config.New(u, config.Token(t, ""), config.SkipTLSValidation())
	} else {
		cfg, err = config.New(u, config.Token(t, ""))
	}
	if err != nil {
		return nil, err
	}

	return client.New(cfg)
}