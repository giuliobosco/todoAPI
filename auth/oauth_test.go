package auth

import (
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/giuliobosco/todoAPI/tu"
	"github.com/stretchr/testify/assert"
)

// ################# TESTS
// SetupOAuth()

// TestSetupOAuthNoCredsFilePath tests SetupOAuth() func without the creds file path
func TestSetupOAuthNoCredsFilePath(t *testing.T) {
	// Running the method in a SubProcess
	if os.Getenv("BE_CRASHER") == "1" {
		SetupOAuth("")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSetupOAuth")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && e.Success() {
		return
	}

	assert.Equal(t, "exit status 1", err.Error())
}

// TestSetupOAuth test SetupOAuth() func with creds file path
func TestSetupOAuth(t *testing.T) {
	SetupOAuth("../creds.json")

	assert.NotNil(t, cred)
	assert.NotNil(t, conf)
}

// ################# TESTS
// SetupOAuth()

// TestOAuthURL tests OAuthURL() returns
func TestOAuthURL(t *testing.T) {
	SetupOAuth("../creds.json")

	w, c := tu.GetRecorderContext()

	req, _ := http.NewRequest("GET", "/", nil)

	c.Request = req

	OAuthURL(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "https://accounts.google.com")
}

// ################# TESTS
// SetupOAuth()

// TestOauthAuthenticator tests OAuthAuthenticator without code param
func TestOauthAuthenticatorNoCodeParam(t *testing.T) {
	SetupOAuth("../creds.json")
	c := tu.GetContext()

	req, _ := http.NewRequest("GET", "/", nil)

	c.Request = req

	u, err := OAuthAuthenticator(c)

	assert.Nil(t, u)
	assert.NotNil(t, err)
}

// TODO other tests
