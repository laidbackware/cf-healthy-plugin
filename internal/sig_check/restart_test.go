package sig_check

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"testing"

// 	"github.com/cloudfoundry/go-cfclient/v3/client"
// 	"github.com/cloudfoundry/go-cfclient/v3/config"
// 	"github.com/stretchr/testify/assert"
// )

// func initClient(t *testing.T) *client.Client {
// 	api_endpoint := mustEnv(t, "CF_API")
// 	username := mustEnv(t, "CF_USER")
// 	password := mustEnv(t, "CF_PASS")

// 	cfg, _ := config.New(fmt.Sprintf("https://%s", api_endpoint),
// 		config.UserPassword(username, password),
// 		config.SkipTLSValidation(),
// 	)
// 	cf, _ := client.New(cfg)
// 	return cf
// }

// func mustEnv(t *testing.T, k string) string {
// 	if v, ok := os.LookupEnv(k); ok {
// 		return v
// 	}
// 	t.Fatalf("expected environment variable %q", k)
// 	return ""
// }

// func TestRestartApp(t *testing.T) {
// 	cf := initClient(t)
// 	l := log.New(os.Stderr, "", 0)
// 	err := restartApp(cf, context.Background(), "9798e4a9-b118-4bd7-ab5f-ac9c6f839b65", l)
// 	assert.Nil(t, err)
// }
