package collect_data

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/config"
	"github.com/stretchr/testify/assert"
)

func initClient(t *testing.T) *client.Client {
	api_endpoint := mustEnv(t, "CF_API")
	username := mustEnv(t, "CF_USER")
	password := mustEnv(t, "CF_PASS")

	cfg, _ := config.New(fmt.Sprintf("https://%s", api_endpoint),
		config.UserPassword(username, password),
		config.SkipTLSValidation(),
	)
	cf, _ := client.New(cfg)
	return cf
}

func mustEnv(t *testing.T, k string) string {
	// t.Helper()
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	t.Fatalf("expected environment variable %q", k)
	return ""
}

func TestAppNameLookup(t *testing.T) {
	cf := initClient(t)
	orgLookup, err := OrgNameLookup(cf)
	assert.Nil(t, err)
	spaceLookup, spaceOrgNameLookup, err := SpaceLookup(cf, orgLookup)
	assert.Nil(t, err)
	appLookup, appSpaceNameLookup, appOrgNameLookup, err := AppLookup(cf, spaceOrgNameLookup, spaceLookup)
	assert.Nil(t, err)
	assert.Greater(t, len(appLookup), 2)
	assert.Greater(t, len(appSpaceNameLookup), 2)
	assert.Greater(t, len(appOrgNameLookup), 2)
}
