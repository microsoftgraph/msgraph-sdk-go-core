package msgraphgocore

import (
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

type batchResponse struct {
	responses     []BatchItem
	indexResponse map[string]BatchItem
	isIndexed     bool
}

func (br *batchResponse) GetResponses() []BatchItem {
	return br.responses
}

func (br *batchResponse) SetResponses(responses []BatchItem) {
	br.responses = responses
}

func (br *batchResponse) GetResponseById(itemId string) BatchItem {
	if !br.isIndexed {

		for _, resp := range br.GetResponses() {
			br.indexResponse[*(resp.GetId())] = resp
		}

		br.isIndexed = true
	}

	return br.indexResponse[itemId]
}

func CreateBatchResponseDiscriminator(serialization.ParseNode) (serialization.Parsable, error) {
	res := batchResponse{
		indexResponse: make(map[string]BatchItem),
		isIndexed:     false,
	}
	return &res, nil
}

type BatchResponse interface {
	serialization.Parsable
	GetResponses() []BatchItem
	SetResponses(responses []BatchItem)
	GetResponseById(itemId string) BatchItem
}

func (br *batchResponse) Serialize(writer serialization.SerializationWriter) error {
	if br.GetResponses() != nil {
		cast := make([]serialization.Parsable, len(br.GetResponses()))
		for i, v := range br.GetResponses() {
			cast[i] = v.(serialization.Parsable)
		}
		err := writer.WriteCollectionOfObjectValues("responses", cast)
		if err != nil {
			return err
		}
	}
	return nil
}

func (br *batchResponse) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	res := make(map[string]func(serialization.ParseNode) error)
	res["responses"] = func(n serialization.ParseNode) error {
		val, err := n.GetCollectionOfObjectValues(CreateBatchRequestItemDiscriminator)
		if err != nil {
			return err
		}
		if val != nil {
			res := make([]BatchItem, len(val))
			for i, v := range val {
				res[i] = v.(BatchItem)
			}
			br.SetResponses(res)
		}
		return nil
	}
	return res
}
