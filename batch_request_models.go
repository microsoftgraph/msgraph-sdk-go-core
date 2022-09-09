package msgraphgocore

import (
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

type Header map[string]string

func (br Header) Serialize(writer serialization.SerializationWriter) error {
	if br != nil {
		for key, element := range br {
			err := writer.WriteStringValue(key, &element)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (br Header) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	res := make(map[string]func(serialization.ParseNode) error)
	if br != nil {
		for key, element := range br {
			res[key] = func(n serialization.ParseNode) error {
				return SetValue(n.GetStringValue, func(val *string) {
					br[key] = element
				})
			}
		}
	}
	return res
}

type RequestBody map[string]interface{}

func (br RequestBody) Serialize(writer serialization.SerializationWriter) error {
	if br != nil {
		err := writer.WriteAdditionalData(br)
		if err != nil {
			return err
		}
	}
	return nil
}

func (br RequestBody) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	res := make(map[string]func(serialization.ParseNode) error)
	if br != nil {
		for key, element := range br {
			res[key] = func(n serialization.ParseNode) error {
				return SetValue(n.GetStringValue, func(val *string) {
					br[key] = element
				})
			}
		}
	}
	return res
}

type BatchItem interface {
	GetId() *string
	SetId(value *string)
	GetMethod() *string
	SetMethod(value *string)
	GetUrl() *string
	SetUrl(value *string)
	GetHeaders() Header
	SetHeaders(value Header)
	GetBody() RequestBody
	SetBody(value RequestBody)
	GetDependsOn() []string
	SetDependsOn(value []string)
	GetStatus() *int32
	SetStatus(value *int32)
	DependsOnItem(item BatchItem)
}

type batchItem struct {
	Id        *string
	method    *string
	Url       *string
	Headers   Header
	Body      RequestBody
	DependsOn []string
	Status    *int32
}

func (bi *batchItem) GetId() *string {
	return bi.Id
}

func (bi *batchItem) SetId(value *string) {
	bi.Id = value
}

func (bi *batchItem) GetMethod() *string {
	return bi.method
}

func (bi *batchItem) SetMethod(value *string) {
	bi.method = value
}

func (bi *batchItem) GetUrl() *string {
	return bi.Url
}

func (bi *batchItem) SetUrl(value *string) {
	bi.Url = value
}

func (bi *batchItem) GetHeaders() Header {
	return bi.Headers
}

func (bi *batchItem) SetHeaders(value Header) {
	bi.Headers = value
}

func (bi *batchItem) GetBody() RequestBody {
	return bi.Body
}

func (bi *batchItem) SetBody(value RequestBody) {
	bi.Body = value
}

func (bi *batchItem) GetDependsOn() []string {
	return bi.DependsOn
}

func (bi *batchItem) SetDependsOn(value []string) {
	bi.DependsOn = value
}

func (bi *batchItem) GetStatus() *int32 {
	return bi.Status
}

func (bi *batchItem) SetStatus(value *int32) {
	bi.Status = value
}

func CreateHeaderFromDiscriminatorValue(serialization.ParseNode) (serialization.Parsable, error) {
	var res Header = make(map[string]string)
	return res, nil
}

func CreateRequestBodyFromDiscriminatorValue(node serialization.ParseNode) (serialization.Parsable, error) {
	var res RequestBody = make(map[string]interface{})
	return res, nil
}

func (bi *batchItem) Serialize(writer serialization.SerializationWriter) error {
	{
		err := writer.WriteStringValue("id", bi.GetId())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteStringValue("method", bi.GetMethod())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteStringValue("url", bi.GetUrl())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteObjectValue("headers", bi.GetHeaders())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteObjectValue("body", bi.GetBody())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteCollectionOfStringValues("dependsOn", bi.GetDependsOn())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteInt32Value("status", bi.GetStatus())
		if err != nil {
			return err
		}
	}
	return nil
}

func (bi *batchItem) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	res := make(map[string]func(serialization.ParseNode) error)
	res["id"] = func(n serialization.ParseNode) error { return SetValue(n.GetStringValue, bi.SetId) }
	res["method"] = func(n serialization.ParseNode) error { return SetValue(n.GetStringValue, bi.SetMethod) }
	res["url"] = func(n serialization.ParseNode) error { return SetValue(n.GetStringValue, bi.SetUrl) }
	res["headers"] = func(n serialization.ParseNode) error {
		return SetObjectValue(n.GetObjectValue, CreateHeaderFromDiscriminatorValue, bi.SetHeaders)
	}
	res["body"] = func(n serialization.ParseNode) error {
		return SetObjectValue(n.GetObjectValue, CreateRequestBodyFromDiscriminatorValue, bi.SetBody)
	}
	//res["dependsOn"] = func(n serialization.ParseNode) error {
	////return SetValue(n.GetCollectionOfPrimitiveValues(), bi.SetDependsOn)
	//}
	res["status"] = func(n serialization.ParseNode) error { return SetValue(n.GetInt32Value, bi.SetStatus) }
	return res
}

func CreateBatchRequestItemDiscriminator(serialization.ParseNode) (serialization.Parsable, error) {
	var res batchItem
	return &res, nil
}

type batchRequest struct {
	Requests []BatchItem
}

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
