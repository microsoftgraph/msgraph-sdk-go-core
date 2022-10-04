package msgraphgocore

import (
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegistration(t *testing.T) {
	errorMapping := abstractions.ErrorMappings{}
	err := RegisterError(BatchRequestErrorRegistryKey, errorMapping)
	require.NoError(t, err)
	err = RegisterError(BatchRequestErrorRegistryKey, errorMapping)
	assert.Equal(t, err.Error(), "object Factory already register")

	err = DeRegisterError(BatchRequestErrorRegistryKey)
	require.NoError(t, err)
	err = DeRegisterError(BatchRequestErrorRegistryKey)
	assert.Equal(t, err.Error(), "object Factory does not exist register")
}
