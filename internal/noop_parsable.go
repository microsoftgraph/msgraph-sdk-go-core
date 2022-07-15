package internal

import (
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

// NoOpParsable is a dummy type that implements the Parsable interface.
// Used in SendBatch to by-pass serialization
type NoOpParsable struct {
}

// NewNoOpParsable creates an instance of NoOpParsable
func NewNoOpParsable() *NoOpParsable {
	return &NoOpParsable{}
}

func (r *NoOpParsable) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	res := make(map[string]func(serialization.ParseNode) error)
	return res
}

func (r *NoOpParsable) Serialize(writer serialization.SerializationWriter) error {
	return nil
}
