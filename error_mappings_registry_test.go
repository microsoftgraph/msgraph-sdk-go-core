package msgraphgocore

import (
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegistration(t *testing.T) {
	errorMapping := abstractions.ErrorMappings{}
	err := RegisterError(BATCH_REQUEST_ERROR_REGISTRY_KEY, errorMapping)
	require.NoError(t, err)
	err = RegisterError(BATCH_REQUEST_ERROR_REGISTRY_KEY, errorMapping)
	assert.Equal(t, err.Error(), "object Factory already register")

	err = DeRegisterError(BATCH_REQUEST_ERROR_REGISTRY_KEY)
	require.NoError(t, err)
	err = DeRegisterError(BATCH_REQUEST_ERROR_REGISTRY_KEY)
	assert.Equal(t, err.Error(), "object Factory does not exist register")
}
