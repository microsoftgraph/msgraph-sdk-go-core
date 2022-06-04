package msgraphgocore

import (
	"fmt"

	"github.com/microsoft/kiota-abstractions-go/serialization"
)

type BatchResponse struct {
	Responses []BatchItem
}

type BatchItemResponse struct {
	Id      *string
	Status  *int32
	Body    *string
	Headers *string
}

func NewBatchResponse() *BatchResponse {
	return &BatchResponse{make([]BatchItem, 0)}
}

func (r *BatchItemResponse) GetBody() *string {
	return r.Body
}

func (r *BatchItemResponse) GetHeaders() *string {
	return r.Headers
}

func (r *BatchItemResponse) GetId() *string {
	return r.Id
}

func (r *BatchItemResponse) GetStatus() *int32 {
	return r.Status
}

func (r *BatchItemResponse) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	res := make(map[string]func(serialization.ParseNode) error)
	res["id"] = func(n serialization.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			r.Id = val
		}
		return nil
	}

	res["status"] = func(n serialization.ParseNode) error {
		val, err := n.GetInt32Value()
		if err != nil {
			return err
		}
		if val != nil {
			r.Status = val
		}
		return nil
	}

	res["body"] = func(n serialization.ParseNode) error {
		val, _ := n.GetStringValue()
		fmt.Println(string(*val))
		fmt.Println("-----")
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			r.Body = val
		}
		return nil
	}

	res["headers"] = func(n serialization.ParseNode) error {
		// val, err := n.GetStringValue()
		// if err != nil {
		// 	return err
		// }
		// if val != nil {
		// 	r.Headers = val
		// }
		return nil
	}
	return res
}

func (r *BatchItemResponse) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

func CreateBatchItem(constructor serialization.ParseNode) (serialization.Parsable, error) {
	return &BatchItemResponse{}, nil
}

func (r *BatchResponse) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	res := make(map[string]func(serialization.ParseNode) error)
	res["responses"] = func(n serialization.ParseNode) error {
		val, err := n.GetCollectionOfObjectValues(CreateBatchItem)
		if err != nil {
			return err
		}
		if val != nil {
			items := make([]BatchItem, len(val))
			for i, item := range val {
				items[i] = item.(BatchItem)
			}

			r.Responses = items
		}
		return nil
	}

	return res
}

func (r *BatchResponse) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

type BatchItem interface {
	GetBody() *string
	GetStatus() *int32
	GetHeaders() *string
	GetId() *string
}
