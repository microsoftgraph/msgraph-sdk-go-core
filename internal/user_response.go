package internal

import (
	i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UsersResponse
type UsersResponse struct {
	// Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
	additionalData map[string]interface{}
	//
	odataNextLink *string
	//
	value []User
}

// NewUsersResponse instantiates a new usersResponse and sets the default values.
func NewUsersResponse() *UsersResponse {
	m := &UsersResponse{}
	m.SetAdditionalData(make(map[string]interface{}))
	return m
}

// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *UsersResponse) GetAdditionalData() map[string]interface{} {
	if m == nil {
		return nil
	} else {
		return m.additionalData
	}
}

// GetOdataNextLink gets the @odata.nextLink property value.
func (m *UsersResponse) GetOdataNextLink() *string {
	if m == nil {
		return nil
	} else {
		return m.odataNextLink
	}
}

// GetValue gets the value property value.
func (m *UsersResponse) GetValue() []User {
	if m == nil {
		return nil
	} else {
		return m.value
	}
}

// GetFieldDeserializers the deserialization information for the current model
func (m *UsersResponse) GetFieldDeserializers() map[string]func(i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) error {
	res := make(map[string]func(i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) error)
	res["@odata.nextLink"] = func(n i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetOdataNextLink(val)
		}
		return nil
	}
	res["value"] = func(n i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) error {
		val, err := n.GetCollectionOfObjectValues(func(pn i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) (i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.Parsable, error) {
			return NewUser(), nil
		})
		if err != nil {
			return err
		}
		if val != nil {
			res := make([]User, len(val))
			for i, v := range val {
				res[i] = *(v.(*User))
			}
			m.SetValue(res)
		}
		return nil
	}
	return res
}
func (m *UsersResponse) IsNil() bool {
	return m == nil
}

// Serialize serializes information the current object
func (m *UsersResponse) Serialize(writer i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.SerializationWriter) error {
	{
		err := writer.WriteStringValue("@odata.nextLink", m.GetOdataNextLink())
		if err != nil {
			return err
		}
	}
	{
		cast := make([]i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.Parsable, len(m.GetValue()))
		for i, v := range m.GetValue() {
			temp := v
			cast[i] = i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.Parsable(&temp)
		}
		err := writer.WriteCollectionOfObjectValues("value", cast)
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
	return nil
}

// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *UsersResponse) SetAdditionalData(value map[string]interface{}) {
	if m != nil {
		m.additionalData = value
	}
}

// SetOdataNextLink sets the @odata.nextLink property value.
func (m *UsersResponse) SetOdataNextLink(value *string) {
	if m != nil {
		m.odataNextLink = value
	}
}

// SetValue sets the value property value.
func (m *UsersResponse) SetValue(value []User) {
	if m != nil {
		m.value = value
	}
}
