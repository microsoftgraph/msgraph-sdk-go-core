package msgraphgocore

import (
	"errors"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-sdk-go-core/internal"
	"github.com/stretchr/testify/assert"
	"strconv"
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

func createCallRecordNode(parseNode serialization.ParseNode) (serialization.Parsable, error) {
	return internal.NewCallRecord(), nil
}

func getObjectValue(ctor serialization.ParsableFactory) (serialization.Parsable, error) {
	return internal.NewCallRecord(), nil
}

func TestSetObjectValueWithoutError(t *testing.T) {

	person := internal.NewPerson()
	err := SetObjectValue(getObjectValue, createCallRecordNode, person.SetCallRecord)
	assert.Nil(t, err)
	assert.NotNil(t, person.GetCallRecord())
}

func getObjectsValues(ctor serialization.ParsableFactory) ([]serialization.Parsable, error) {
	slice := []serialization.Parsable{internal.NewCallRecord(), internal.NewCallRecord(), internal.NewCallRecord()}
	return slice, nil
}

func TestSetCollectionValueValueWithoutError(t *testing.T) {

	person := internal.NewPerson()
	err := SetCollectionValue(getObjectsValues, createCallRecordNode, person.SetCallRecords)
	assert.Nil(t, err)
	assert.Equal(t, len(person.GetCallRecords()), 3)
}

func TestCollectionApply(t *testing.T) {

	slice := []string{"1", "2", "3"}
	response := CollectionApply(slice, func(s string) int {
		i, _ := strconv.Atoi(s)
		return i
	})

	assert.Equal(t, len(response), 3)
	assert.Equal(t, response, []int{1, 2, 3})
}
