package internal

import (
	i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UsersDeltaResponse
type UsersDeltaResponse struct {
	//
	UsersResponse
	//
	odataDeltaLink *string
}

// NewDeltasResponse instantiates a new usersResponse and sets the default values.
func NewUsersDeltaResponse() *UsersDeltaResponse {
	m := &UsersDeltaResponse{
		UsersResponse: *NewUsersResponse(),
	}
	return m
}

// GetOdataDeltaLink gets the @odata.nextLink property value.
func (m *UsersDeltaResponse) GetOdataDeltaLink() *string {
	if m == nil {
		return nil
	} else {
		return m.odataDeltaLink
	}
}

// GetFieldDeserializers the deserialization information for the current model
func (m *UsersDeltaResponse) GetFieldDeserializers() map[string]func(i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) error {
	res := m.UsersResponse.GetFieldDeserializers()
	res["@odata.deltaLink"] = func(n i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetOdataDeltaLink(val)
		}
		return nil
	}
	return res
}

// Serialize serializes information the current object
func (m *UsersDeltaResponse) Serialize(writer i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.SerializationWriter) error {
	{
		err := writer.WriteStringValue("@odata.deltaLink", m.GetOdataDeltaLink())
		if err != nil {
			return err
		}
	}

	return m.UsersResponse.Serialize(writer)
}

// SetOdataDeltaLink sets the @odata.deltaLink property value.
func (m *UsersDeltaResponse) SetOdataDeltaLink(value *string) {
	if m != nil {
		m.odataDeltaLink = value
	}
}
