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

func getObjectValueWithError(ctor serialization.ParsableFactory) (serialization.Parsable, error) {
	return nil, errors.New("could not get from factory")
}

func TestSetObjectValueWithoutError(t *testing.T) {

	person := internal.NewPerson()
	err := SetObjectValue(getObjectValue, createCallRecordNode, person.SetCallRecord)
	assert.Nil(t, err)
	assert.NotNil(t, person.GetCallRecord())
}

func TestSetObjectValueWithError(t *testing.T) {

	person := internal.NewPerson()
	err := SetObjectValue(getObjectValueWithError, createCallRecordNode, person.SetCallRecord)

	assert.NotNil(t, err)
	assert.Nil(t, person.GetCallRecord())
}

func getObjectsValues(ctor serialization.ParsableFactory) ([]serialization.Parsable, error) {
	slice := []serialization.Parsable{internal.NewCallRecord(), internal.NewCallRecord(), internal.NewCallRecord()}
	return slice, nil
}

func getObjectsValuesWithError(ctor serialization.ParsableFactory) ([]serialization.Parsable, error) {
	return nil, errors.New("could not get from factory")
}

func TestSetCollectionValueValueWithoutError(t *testing.T) {

	person := internal.NewPerson()
	err := SetCollectionValue(getObjectsValues, createCallRecordNode, person.SetCallRecords)
	assert.Nil(t, err)
	assert.NotNil(t, person.GetCallRecords())
	assert.Equal(t, len(person.GetCallRecords()), 3)
}

func TestSetCollectionValueValueWithError(t *testing.T) {

	person := internal.NewPerson()
	err := SetCollectionValue(getObjectsValuesWithError, createCallRecordNode, person.SetCallRecords)
	assert.NotNil(t, err)
	assert.Nil(t, person.GetCallRecords())
	assert.Equal(t, len(person.GetCallRecords()), 0)
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

func getEnumValue(parser serialization.EnumFactory) (interface{}, error) {
	status := internal.ACTIVE
	return &status, nil
}

func TestSetEnumValueValueValueWithoutError(t *testing.T) {
	person := internal.NewPerson()
	err := SetEnumValue(getEnumValue, internal.ParsePersonStatus, person.SetStatus)
	assert.Nil(t, err)
	assert.Equal(t, person.GetStatus().String(), internal.ACTIVE.String())
}

func getEnumValueWithError(parser serialization.EnumFactory) (interface{}, error) {
	return nil, errors.New("could not get from factory")
}

func TestSetEnumValueValueValueWithError(t *testing.T) {
	person := internal.NewPerson()
	err := SetEnumValue(getEnumValueWithError, internal.ParsePersonStatus, person.SetStatus)
	assert.NotNil(t, err)
	assert.Nil(t, person.GetStatus())
}

func TestSetReferencedEnumValueValueValueWithoutError(t *testing.T) {
	person := internal.NewPerson()
	err := SetReferencedEnumValue(getEnumValue, internal.ParsePersonStatus, person.SetStatus)
	assert.Nil(t, err)
	assert.Equal(t, person.GetStatus().String(), internal.ACTIVE.String())
}

func TestSetReferencedEnumValueValueValueWithError(t *testing.T) {
	person := internal.NewPerson()
	err := SetEnumValue(getEnumValueWithError, internal.ParsePersonStatus, person.SetStatus)
	assert.NotNil(t, err)
	assert.Nil(t, person.GetStatus())
}

func TestSetCollectionOfReferencedEnumValueWithoutError(t *testing.T) {
	person := internal.NewPerson()

	slice := []interface{}{Point(internal.ACTIVE), Point(internal.SUSPEND)}
	enumSource := func(parser serialization.EnumFactory) ([]interface{}, error) {
		return slice, nil
	}

	err := SetCollectionOfReferencedEnumValue(enumSource, internal.ParsePersonStatus, person.SetPreviousStatus)
	assert.Nil(t, err)
	assert.Equal(t, person.GetPreviousStatus()[0].String(), internal.ACTIVE.String())
	assert.Equal(t, person.GetPreviousStatus()[1].String(), internal.SUSPEND.String())
}

func TestSetCollectionOfReferencedEnumValueWithError(t *testing.T) {
	person := internal.NewPerson()

	slice := []interface{}{Point(internal.ACTIVE), Point(internal.SUSPEND)}
	enumSource := func(parser serialization.EnumFactory) ([]interface{}, error) {
		return slice, nil
	}

	err := SetCollectionOfReferencedEnumValue(enumSource, internal.ParsePersonStatus, person.SetPreviousStatus)
	assert.NotNil(t, err)
	assert.Nil(t, person.GetPreviousStatus())
}

func TestSetSetCollectionOfPrimitiveValueWithoutError(t *testing.T) {
	person := internal.NewPerson()

	slice := []interface{}{Point(1), Point(2), Point(3)}
	dataSource := func(targetType string) ([]interface{}, error) {
		return slice, nil
	}

	err := SetCollectionOfReferencedPrimitiveValue(dataSource, "int", person.SetCardNumbers)
	assert.Nil(t, err)
	assert.Equal(t, person.GetCardNumbers()[0], 1)
	assert.Equal(t, person.GetCardNumbers()[1], 2)
	assert.Equal(t, person.GetCardNumbers()[2], 3)
}

func TestSetSetCollectionOfPrimitiveValueWithError(t *testing.T) {
	person := internal.NewPerson()

	dataSource := func(targetType string) ([]interface{}, error) {
		return nil, errors.New("could not get from factory")
	}

	err := SetCollectionOfReferencedPrimitiveValue(dataSource, "int", person.SetCardNumbers)
	assert.NotNil(t, err)
	assert.Nil(t, person.GetCardNumbers())
}
