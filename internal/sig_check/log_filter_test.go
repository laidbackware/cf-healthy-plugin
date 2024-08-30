package sig_check

// import (
// 	"crypto/tls"
// 	"net/http"
// 	"sync"
// 	"testing"
// 	"time"

// 	"code.cloudfoundry.org/cli/plugin"
// 	"github.com/stretchr/testify/assert"
// )

// type stubCliConnection struct {
// 	plugin.CliConnection
// 	sync.Mutex

// 	apiEndpointErr error
// 	apiEndpoint		 string

// 	hasAPIEndpoint    bool
// 	hasAPIEndpointErr error

// 	cliCommandArgs   [][]string
// 	cliCommandResult [][]string
// 	cliCommandErr    []error

// 	orgName      string
// 	spaceName    string

// 	accessTokenCount int
// 	accessToken      string
// 	accessTokenErr   error

// 	sslDisabled	bool
// }

// func (s *stubCliConnection) AccessToken() (string, error) {
// 	s.Lock()
// 	defer s.Unlock()

// 	s.accessTokenCount++
// 	return s.accessToken, s.accessTokenErr
// }

// func (s *stubCliConnection) IsSSLDisabled() (bool, error) {
// 	return s.sslDisabled, nil
// }

// func (s *stubCliConnection) ApiEndpoint() (string, error) {
// 	return s.apiEndpoint, s.apiEndpointErr
// }

// func (s *stubCliConnection) HasAPIEndpoint() (bool, error) {
// 	return s.hasAPIEndpoint, s.hasAPIEndpointErr
// }

// func newStubCliConnection() *stubCliConnection {
// 	return &stubCliConnection{
// 		hasAPIEndpoint: true,
// 	}
// }

// func TestGetLogs(t *testing.T) {
// 	cliConnection := newStubCliConnection()
// 	cliConnection.accessToken = "bearer eyJqa3UiOiJodHRwczovL3VhYS5zeXMuMTkyLjE2OC4xLjIyOS5uaXAuaW8vdG9rZW5fa2V5cyIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIwMmJkMjJlNy0zZDBmLTRkOWMtYjUyYy03MDIzMjgxNWM3ZmUiLCJ1c2VyX25hbWUiOiJhZG1pbiIsIm9yaWdpbiI6InVhYSIsImlzcyI6Imh0dHBzOi8vdWFhLnN5cy4xOTIuMTY4LjEuMjI5Lm5pcC5pby9vYXV0aC90b2tlbiIsImNsaWVudF9pZCI6ImNmIiwiYXVkIjpbImNmIiwibmV0d29yayIsInVhYSIsInNjaW0iLCJwYXNzd29yZCIsImNsb3VkX2NvbnRyb2xsZXIiLCJvcGVuaWQiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMiLCJkb3BwbGVyIl0sInppZCI6InVhYSIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiIwMmJkMjJlNy0zZDBmLTRkOWMtYjUyYy03MDIzMjgxNWM3ZmUiLCJhenAiOiJjZiIsInNjb3BlIjpbIm9wZW5pZCIsInJvdXRpbmcucm91dGVyX2dyb3Vwcy53cml0ZSIsIm5ldHdvcmsud3JpdGUiLCJzY2ltLnJlYWQiLCJjbG91ZF9jb250cm9sbGVyLmFkbWluIiwidWFhLnVzZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMucmVhZCIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSIsIm5ldHdvcmsuYWRtaW4iLCJkb3BwbGVyLmZpcmVob3NlIiwic2NpbS53cml0ZSJdLCJhdXRoX3RpbWUiOjE3MjQ4NDk0MzQsImV4cCI6MTcyNDk0MjU5MSwiaWF0IjoxNzI0OTM1MzkxLCJqdGkiOiI1NzU3Y2E0YzdlNWE0MGE1YTM4MmU0ZGE1ZjM1Y2RiMyIsImVtYWlsIjoiYWRtaW4iLCJyZXZfc2lnIjoiZDUxZTNmZTAiLCJjbGllbnRfYXV0aF9tZXRob2QiOiJub25lIiwiY2lkIjoiY2YifQ.GrBi4DO4Cfi4ZQtzghua48Gag_L9fSFYFQkjRV8SZ5nPcspzBb6PnS7ILyiEUuX7GUwpHwOlQNC6ENLmHiAV2tanvBBiXufPgFk5Hul4Pj1Lp0G92d_CTYx_5BIsx9G9of1VfZG4EHcYrxE06g4xj5-OlJb3L8rPdos-Zt4f06KxafLVDpTYBvaOrAArNXTfCNY_jwrRnwv69Sru2NqnX0m37lwa-jPkPRAtYgqCfMSEdAbjCqk2ca4Yfp6hIU-HcaMUihOCbAtrVDXIh1QT_HP9fFN9v3T54EEO69k7lTwyiDvoLeEw4TCG-aYgddGiRrY2ZBJXnhqoYx_pqR0vDg"
// 	cliConnection.hasAPIEndpoint = true
// 	cliConnection.apiEndpoint = "https://api.sys.192.168.1.229.nip.io"
// 	cliConnection.sslDisabled = true
// 	sourceID := "9798e4a9-b118-4bd7-ab5f-ac9c6f839b65"

// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
// 		InsecureSkipVerify: true, //nolint:gosec
// 	}

// 	quit := make(chan bool)
// 	// quit <- false
// 	shutdown := make(chan int)
// 	sigkill := make(chan int)
// 	err := make(chan error)

// 	GetLogs(cliConnection, http.DefaultClient, sourceID, quit, shutdown, sigkill, err, true)

// 	time.Sleep(3 * time.Second)

// 	quit <- true

// 	assert.Nil(t, <-err)
// }
