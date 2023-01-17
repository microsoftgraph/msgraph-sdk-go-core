package authentication

import (
	"testing"

	absauth "github.com/microsoft/kiota-abstractions-go/authentication"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticationProviderImplementsInterface(t *testing.T) {
	var value absauth.AuthenticationProvider = &AzureIdentityAuthenticationProvider{}
	assert.NotNil(t, value)
}
