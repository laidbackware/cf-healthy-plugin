package sig_check

import (
	"crypto/tls"
	"net/http"
	"sync"
	"testing"
	"time"

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
	cliConnection.accessToken = "bearer eyJqa3UiOiJodHRwczovL3VhYS5zeXMuMTkyLjE2OC4xLjIyOS5uaXAuaW8vdG9rZW5fa2V5cyIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIwMmJkMjJlNy0zZDBmLTRkOWMtYjUyYy03MDIzMjgxNWM3ZmUiLCJ1c2VyX25hbWUiOiJhZG1pbiIsIm9yaWdpbiI6InVhYSIsImlzcyI6Imh0dHBzOi8vdWFhLnN5cy4xOTIuMTY4LjEuMjI5Lm5pcC5pby9vYXV0aC90b2tlbiIsImNsaWVudF9pZCI6ImNmIiwiYXVkIjpbImRvcHBsZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMiLCJvcGVuaWQiLCJjbG91ZF9jb250cm9sbGVyIiwicGFzc3dvcmQiLCJzY2ltIiwidWFhIiwibmV0d29yayIsImNmIl0sInppZCI6InVhYSIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiIwMmJkMjJlNy0zZDBmLTRkOWMtYjUyYy03MDIzMjgxNWM3ZmUiLCJhenAiOiJjZiIsInNjb3BlIjpbIm9wZW5pZCIsInJvdXRpbmcucm91dGVyX2dyb3Vwcy53cml0ZSIsIm5ldHdvcmsud3JpdGUiLCJzY2ltLnJlYWQiLCJjbG91ZF9jb250cm9sbGVyLmFkbWluIiwidWFhLnVzZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMucmVhZCIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSIsIm5ldHdvcmsuYWRtaW4iLCJkb3BwbGVyLmZpcmVob3NlIiwic2NpbS53cml0ZSJdLCJhdXRoX3RpbWUiOjE3MjIzMjc3NzksImV4cCI6MTcyMzIxNDg2OSwiaWF0IjoxNzIzMjA3NjY5LCJqdGkiOiI2YmI3NTg4YmY0MzU0NjE2ODUwMzc1MjIzOGJiNWFiYyIsImVtYWlsIjoiYWRtaW4iLCJyZXZfc2lnIjoiZDUxZTNmZTAiLCJjbGllbnRfYXV0aF9tZXRob2QiOiJub25lIiwiY2lkIjoiY2YifQ.Zp7Ck_t8yQ9OAqgrAW2TDk_GGiXvVtGfUv38Z0J_78WbirioLWpmfdXQVeaZppKIs2vQ053dPzkT63BZr5wv6R2aG8gGX9LjD-1LFw4dj1ALLQXp4p5JB4gMoxd_xj8uovxMvUkMLFbEMUoCL1_wQmXthU6i6_pnFsu7XYQowsY56OgR724rbL6YjYjSJyo-zYD_uLT3kJGjxnekKhexwsHCSB1gdQw3rsr6mRM0uJP3lC1n0eZNb8jBPLafNTT0Hz4xv6UMb4b2qLAaBKIXKG4GvMAJK6lrEhaO6RCqmp_Ji9sZkh7GkKGpkd4BYbPuY1M2vmNKWiiL4lnLj6Ct1g"
	cliConnection.hasAPIEndpoint = true
	cliConnection.apiEndpoint = "https://api.sys.192.168.1.229.nip.io"
	cliConnection.sslDisabled = true
	sourceID := "9798e4a9-b118-4bd7-ab5f-ac9c6f839b65"

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true, //nolint:gosec
	}

	quit := make(chan bool)
	shutdown := make(chan int)
	sigkill := make(chan int)
	err := make(chan error)

	go GetLogs(cliConnection, http.DefaultClient, sourceID, quit, shutdown, sigkill, err)

	time.Sleep(3 * time.Second)

	// text/plain

	quit <- true

	assert.Nil(t, <-err)
}