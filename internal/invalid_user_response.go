package internal

// InvalidUsersResponse
type InvalidUsersResponse struct {
	// Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
	additionalData map[string]interface{}
	//
	nextLink *string
	//
	value []User
}

// NewInvalidUsersResponse instantiates a new InvalidUsersResponse and sets the default values.
func NewInvalidUsersResponse() *InvalidUsersResponse {
	m := &InvalidUsersResponse{}
	m.SetAdditionalData(make(map[string]interface{}))
	return m
}

// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *InvalidUsersResponse) GetAdditionalData() map[string]interface{} {
	if m == nil {
		return nil
	} else {
		return m.additionalData
	}
}

// GetNextLink gets the @odata.nextLink property value.
func (m *InvalidUsersResponse) GetNextLink() *string {
	if m == nil {
		return nil
	} else {
		return m.nextLink
	}
}

// GetValue gets the value property value.
func (m *InvalidUsersResponse) GetValue() []User {
	if m == nil {
		return nil
	} else {
		return m.value
	}
}

func (m *InvalidUsersResponse) IsNil() bool {
	return m == nil
}

// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *InvalidUsersResponse) SetAdditionalData(value map[string]interface{}) {
	if m != nil {
		m.additionalData = value
	}
}

// SetNextLink sets the @odata.nextLink property value.
func (m *InvalidUsersResponse) SetNextLink(value *string) {
	if m != nil {
		m.nextLink = value
	}
}

// SetValue sets the value property value.
func (m *InvalidUsersResponse) SetValue(value []User) {
	if m != nil {
		m.value = value
	}
}
