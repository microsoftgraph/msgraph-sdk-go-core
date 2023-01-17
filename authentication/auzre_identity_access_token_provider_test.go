package authentication

import (
	"testing"

	absauth "github.com/microsoft/kiota-abstractions-go/authentication"
	"github.com/stretchr/testify/assert"
)

func TestAccessTokenProviderImplementsInterface(t *testing.T) {
	var value absauth.AccessTokenProvider = &AzureIdentityAccessTokenProvider{}
	assert.NotNil(t, value)
}
