package internal

import (
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

type NoOpParsable struct {
}

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
