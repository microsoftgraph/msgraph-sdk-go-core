package msgraphgocore

import (
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
)

type ParseFactory[T serialization.Parsable] func() T

// ParseCollections represents a contract with the GetValue() method
type ParseCollections[T serialization.Parsable] interface {
	serialization.Parsable
	GetValue() []T
	SetValue(value []T)
}

type ParseableCollection[T serialization.Parsable] struct {
	backingStore    store.BackingStore
	constructorFunc serialization.ParsableFactory
}

// NewParseableCollection instantiates a new ParseableCollection and sets the default values.
func NewParseableCollection[T serialization.Parsable](constructorFunc ParseFactory[T]) *ParseableCollection[T] {
	m := &ParseableCollection[T]{
		backingStore: store.BackingStoreFactoryInstance(),
		constructorFunc: func(node serialization.ParseNode) (serialization.Parsable, error) {
			return constructorFunc(), nil
		},
	}
	m.SetAdditionalData(make(map[string]any))
	return m
}

// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(serialization.ParseNode)(error) when successful
func (m *ParseableCollection[T]) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	res := make(map[string]func(serialization.ParseNode) error)
	res["value"] = func(n serialization.ParseNode) error {
		val, err := n.GetCollectionOfObjectValues(m.constructorFunc)
		if err != nil {
			return err
		}
		if val != nil {
			res := make([]T, len(val))
			for i, v := range val {
				if v != nil {
					res[i] = v.(T)
				}
			}
			m.SetValue(res)
		}
		return nil
	}
	res["@odata.nextLink"] = func(n serialization.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetOdataNextLink(val)
		}
		return nil
	}
	return res
}

// GetValue gets the value property value. The value property
// returns a []Userable when successful
func (m *ParseableCollection[T]) GetValue() []T {
	val, err := m.GetBackingStore().Get("value")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.([]T)
	}
	return nil
}

// Serialize serializes information the current object
func (m *ParseableCollection[T]) Serialize(writer serialization.SerializationWriter) error {
	{
		err := writer.WriteStringValue("@odata.nextLink", m.GetOdataNextLink())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteAdditionalData(m.GetAdditionalData())
		if err != nil {
			return err
		}
	}
	if m.GetValue() != nil {
		cast := make([]serialization.Parsable, len(m.GetValue()))
		for i, v := range m.GetValue() {
			cast[i] = any(v).(serialization.Parsable)
		}
		err := writer.WriteCollectionOfObjectValues("value", cast)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ParseableCollection[T]) GetAdditionalData() map[string]any {
	val, err := m.backingStore.Get("additionalData")
	if err != nil {
		panic(err)
	}
	if val == nil {
		var value = make(map[string]any)
		m.SetAdditionalData(value)
	}
	return val.(map[string]any)
}

// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *ParseableCollection[T]) GetBackingStore() store.BackingStore {
	return m.backingStore
}

// SetValue sets the value property value. The value property
func (m *ParseableCollection[T]) SetValue(value []T) {
	err := m.GetBackingStore().Set("value", value)
	if err != nil {
		panic(err)
	}
}

// GetOdataNextLink gets the @odata.nextLink property value. The OdataNextLink property
// returns a *string when successful
func (m *ParseableCollection[T]) GetOdataNextLink() *string {
	val, err := m.GetBackingStore().Get("odataNextLink")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(*string)
	}
	return nil
}

// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ParseableCollection[T]) SetAdditionalData(value map[string]any) {
	err := m.GetBackingStore().Set("additionalData", value)
	if err != nil {
		panic(err)
	}
}

// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ParseableCollection[T]) SetBackingStore(value store.BackingStore) {
	m.backingStore = value
}

// SetOdataCount sets the @odata.count property value. The OdataCount property
func (m *ParseableCollection[T]) SetOdataCount(value *int64) {
	err := m.GetBackingStore().Set("odataCount", value)
	if err != nil {
		panic(err)
	}
}

// SetOdataNextLink sets the @odata.nextLink property value. The OdataNextLink property
func (m *ParseableCollection[T]) SetOdataNextLink(value *string) {
	err := m.GetBackingStore().Set("odataNextLink", value)
	if err != nil {
		panic(err)
	}
}
