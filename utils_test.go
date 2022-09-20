package msgraphgocore

import (
	"errors"
	"github.com/microsoftgraph/msgraph-sdk-go-core/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetValueWithoutError(t *testing.T) {

	person := internal.NewPerson()
	callFactory := func() (*internal.CallRecord, error) {
		return internal.NewCallRecord(), nil
	}

	err := SetValue(callFactory, person.SetCallRecord)
	assert.Nil(t, err)
	assert.NotNil(t, person.GetCallRecord())
}

func TestSetValueWithError(t *testing.T) {

	person := internal.NewPerson()
	callFactory := func() (*internal.CallRecord, error) {
		return nil, errors.New("could not get from factory")
	}

	err := SetValue(callFactory, person.SetCallRecord)
	assert.NotNil(t, err)
	assert.Nil(t, person.GetCallRecord())
}
