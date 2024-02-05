package internal

import (
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

type SampleError struct {
	abstractions.ApiError
	// Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
	additionalData map[string]interface{}
}

func (s SampleError) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

func (s *SampleError) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	res := make(map[string]func(serialization.ParseNode) error)
	res["error"] = func(n serialization.ParseNode) error {
		v, err := n.GetRawValue()
		if err != nil {
			return err
		}
		if vm, ok := v.(map[string]interface{}); ok {
			if msg, ok := vm["message"]; ok && msg != nil {
				s.Message = *msg.(*string)
			}
		}
		return nil
	}
	return res
}

func CreateSampleErrorFromDiscriminatorValue(parseNode serialization.ParseNode) (serialization.Parsable, error) {
	res := SampleError{}
	return &res, nil
}
