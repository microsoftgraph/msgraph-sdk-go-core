package internal

import (
	i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"

	i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55 "github.com/microsoft/kiota/abstractions/go/serialization"
)

type User struct {
	DisplayName *string
	DirectoryObject
}

func (u *User) GetDisplayName() *string {
	return u.DisplayName
}

var displayName = "A User"

func NewUser() *User {
	return &User{
		DisplayName: &displayName,
	}
}

type DirectoryObject struct {
	Entity
	//
	deletedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
}

// NewDirectoryObject instantiates a new directoryObject and sets the default values.
func NewDirectoryObject() *DirectoryObject {
	m := &DirectoryObject{
		Entity: *NewEntity(),
	}
	return m
}

// GetDeletedDateTime gets the deletedDateTime property value.
func (m *DirectoryObject) GetDeletedDateTime() *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time {
	if m == nil {
		return nil
	} else {
		return m.deletedDateTime
	}
}

// GetFieldDeserializers the deserialization information for the current model
func (m *DirectoryObject) GetFieldDeserializers() map[string]func(interface{}, i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) error {
	res := m.Entity.GetFieldDeserializers()
	res["deletedDateTime"] = func(o interface{}, n i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) error {
		val, err := n.GetTimeValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetDeletedDateTime(val)
		}
		return nil
	}
	return res
}
func (m *DirectoryObject) IsNil() bool {
	return m == nil
}

// Serialize serializes information the current object
func (m *DirectoryObject) Serialize(writer i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.SerializationWriter) error {
	err := m.Entity.Serialize(writer)
	if err != nil {
		return err
	}
	{
		err = writer.WriteTimeValue("deletedDateTime", m.GetDeletedDateTime())
		if err != nil {
			return err
		}
	}
	return nil
}

// SetDeletedDateTime sets the deletedDateTime property value.
func (m *DirectoryObject) SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
	if m != nil {
		m.deletedDateTime = value
	}
}

// Entity
type Entity struct {
	// Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
	additionalData map[string]interface{}
	// Read-only.
	id *string
}

// NewEntity instantiates a new entity and sets the default values.
func NewEntity() *Entity {
	m := &Entity{}
	m.SetAdditionalData(make(map[string]interface{}))
	return m
}

// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Entity) GetAdditionalData() map[string]interface{} {
	if m == nil {
		return nil
	} else {
		return m.additionalData
	}
}

// GetId gets the id property value. Read-only.
func (m *Entity) GetId() *string {
	if m == nil {
		return nil
	} else {
		return m.id
	}
}

// GetFieldDeserializers the deserialization information for the current model
func (m *Entity) GetFieldDeserializers() map[string]func(interface{}, i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) error {
	res := make(map[string]func(interface{}, i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) error)
	res["id"] = func(o interface{}, n i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetId(val)
		}
		return nil
	}
	return res
}
func (m *Entity) IsNil() bool {
	return m == nil
}

// Serialize serializes information the current object
func (m *Entity) Serialize(writer i04eb5309aeaafadd28374d79c8471df9b267510b4dc2e3144c378c50f6fd7b55.SerializationWriter) error {
	{
		err := writer.WriteStringValue("id", m.GetId())
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
func (m *Entity) SetAdditionalData(value map[string]interface{}) {
	if m != nil {
		m.additionalData = value
	}
}

// SetId sets the id property value. Read-only.
func (m *Entity) SetId(value *string) {
	if m != nil {
		m.id = value
	}
}
