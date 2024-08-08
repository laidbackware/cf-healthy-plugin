package restart

import (
	"crypto/tls"
	"net/http"
	"sync"
	"testing"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/stretchr/testify/assert"
)

type stubCliConnection struct {
	plugin.CliConnection
	sync.Mutex

	apiEndpointErr error
	apiEndpoint		 string

	hasAPIEndpoint    bool
	hasAPIEndpointErr error

	cliCommandArgs   [][]string
	cliCommandResult [][]string
	cliCommandErr    []error

	orgName      string
	spaceName    string

	accessTokenCount int
	accessToken      string
	accessTokenErr   error

	sslDisabled	bool
}

func (s *stubCliConnection) AccessToken() (string, error) {
	s.Lock()
	defer s.Unlock()

	s.accessTokenCount++
	return s.accessToken, s.accessTokenErr
}

func (s *stubCliConnection) IsSSLDisabled() (bool, error) {
	return s.sslDisabled, nil
}

func (s *stubCliConnection) ApiEndpoint() (string, error) {
	return s.apiEndpoint, s.apiEndpointErr
}

func (s *stubCliConnection) HasAPIEndpoint() (bool, error) {
	return s.hasAPIEndpoint, s.hasAPIEndpointErr
}

func newStubCliConnection() *stubCliConnection {
	return &stubCliConnection{
		hasAPIEndpoint: true,
	}
}

func TestGetLogs(t *testing.T) {
	cliConnection := newStubCliConnection()
	cliConnection.accessToken = "bearer eyJqa3UiOiJodHRwczovL3VhYS5zeXMuMTkyLjE2OC4xLjIyOS5uaXAuaW8vdG9rZW5fa2V5cyIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIwMmJkMjJlNy0zZDBmLTRkOWMtYjUyYy03MDIzMjgxNWM3ZmUiLCJ1c2VyX25hbWUiOiJhZG1pbiIsIm9yaWdpbiI6InVhYSIsImlzcyI6Imh0dHBzOi8vdWFhLnN5cy4xOTIuMTY4LjEuMjI5Lm5pcC5pby9vYXV0aC90b2tlbiIsImNsaWVudF9pZCI6ImNmIiwiYXVkIjpbImRvcHBsZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMiLCJvcGVuaWQiLCJjbG91ZF9jb250cm9sbGVyIiwicGFzc3dvcmQiLCJzY2ltIiwidWFhIiwibmV0d29yayIsImNmIl0sInppZCI6InVhYSIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiIwMmJkMjJlNy0zZDBmLTRkOWMtYjUyYy03MDIzMjgxNWM3ZmUiLCJhenAiOiJjZiIsInNjb3BlIjpbIm9wZW5pZCIsInJvdXRpbmcucm91dGVyX2dyb3Vwcy53cml0ZSIsIm5ldHdvcmsud3JpdGUiLCJzY2ltLnJlYWQiLCJjbG91ZF9jb250cm9sbGVyLmFkbWluIiwidWFhLnVzZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMucmVhZCIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSIsIm5ldHdvcmsuYWRtaW4iLCJkb3BwbGVyLmZpcmVob3NlIiwic2NpbS53cml0ZSJdLCJhdXRoX3RpbWUiOjE3MjIzMjc3NzksImV4cCI6MTcyMzA1MDYzNCwiaWF0IjoxNzIzMDQzNDM0LCJqdGkiOiI0MGVkZDNjYTJjM2Q0ZDcxYmY0OTA3Y2Y4MmQ2YzQxNiIsImVtYWlsIjoiYWRtaW4iLCJyZXZfc2lnIjoiZDUxZTNmZTAiLCJjbGllbnRfYXV0aF9tZXRob2QiOiJub25lIiwiY2lkIjoiY2YifQ.MGxa5lHzm1cEy3dQwVHK5VhjokVE4VXqeO_50VEy6DfgUw8v1Ay8hhQcdMVbE2cnXC_-2qrQEkQD8aR3IqgavICMfsGpSlZVtWD0JqzsPBjnVMA83eEUSfFuQP07sLHlIqoZDz0wvPFgR7g-Zuvn8UuGPiauu9efgj3qUg-TXRYifEaHvRHUGNHc6h5S6piwdOBuyyFaBIr4N9N_CtGo7dbnXtH1kWc9a-TDVMcbaOMTRumAtzF9wtoVKXT9yrjO4oQ1UqrjPFtbUzNmQH-_ZA4nbdmoUA_oG3WkspKtRtgYu6IIwmbGLhS693qUTVvyeAfEyM6gGjk5H0_a0nNAVA"
	cliConnection.hasAPIEndpoint = true
	cliConnection.apiEndpoint = "https://api.sys.192.168.1.229.nip.io"
	cliConnection.sslDisabled = true

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true, //nolint:gosec
	}


	err := GetLogs(cliConnection, http.DefaultClient, "013a044c-4275-4444-8b98-4e70d5842895")
	assert.Nil(t, err)
}