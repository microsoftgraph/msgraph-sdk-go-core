package msgraphgocore

import (
	"errors"
	abs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	"reflect"
)

type BatchItem interface {
	serialization.Parsable
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

// DependsOnItem creates a dependency chain between BatchItems.If A depends on B, then B will be sent before B
// A batchItem can only depend on one other batchItem
// see: https://docs.microsoft.com/en-us/graph/known-issues#request-dependencies-are-limited
func (bi *batchItem) DependsOnItem(item BatchItem) {
	// DependsOn is a errorRegistry value slice
	bi.DependsOn = []string{*item.GetId()}
}

// NewBatchItem creates an instance of BatchItem
func NewBatchItem() BatchItem {
	return &batchItem{
		DependsOn: make([]string, 0),
	}
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
		err := writer.WriteAnyValue("headers", bi.GetHeaders())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteAnyValue("body", bi.GetBody())
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
	res["id"] = abs.SetStringValue(bi.SetId)
	res["method"] = abs.SetStringValue(bi.SetMethod)
	res["url"] = abs.SetStringValue(bi.SetUrl)
	res["headers"] = func(n serialization.ParseNode) error {
		rawVal, err := n.GetRawValue()
		if err != nil {
			return err
		}

		if rawVal == nil {
			return nil
		}

		result, err := castMapOfStrings(rawVal)
		if err != nil {
			return err
		}

		bi.SetHeaders(result)
		return nil
	}
	res["body"] = func(n serialization.ParseNode) error {
		rawVal, err := n.GetRawValue()
		if err != nil {
			return err
		}

		if rawVal == nil {
			return nil
		}

		result, err := convertToMap(rawVal)
		if err != nil {
			return err
		}

		bi.SetBody(result)
		return nil
	}
	res["dependsOn"] = abs.SetCollectionOfPrimitiveValues("string", bi.SetDependsOn)
	res["status"] = abs.SetInt32Value(bi.SetStatus)
	return res
}

func convertToMap(rawVal interface{}) (map[string]interface{}, error) {
	kind := reflect.ValueOf(rawVal)
	if kind.Kind() == reflect.Map {
		result := make(map[string]interface{})
		err := deserializeMapped(kind, result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
	return nil, errors.New("interface was not a map")
}

func deserializeNode(value serialization.ParseNode) (interface{}, error) {
	rawVal, err := value.GetRawValue()
	if err != nil {
		return nil, err
	} else {
		kind := reflect.ValueOf(rawVal)
		if kind.Kind() == reflect.Map {

			result := make(map[string]interface{})
			err := deserializeMapped(kind, result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return deserializeValue(rawVal)
		}
	}
}

func deserializeMapped(v reflect.Value, result map[string]interface{}) error {
	for _, key := range v.MapKeys() {
		value, err := deserializeValue(v.MapIndex(key).Interface())
		if err != nil {
			return err
		} else {
			result[key.String()] = value
		}
	}
	return nil
}

func deserializeNodes(value []*jsonserialization.JsonParseNode) (interface{}, error) {
	slice := make([]interface{}, len(value))
	for index, element := range value {
		res, err := deserializeNode(element)
		if err != nil {
			return nil, err
		}
		slice[index] = res
	}
	return slice, nil
}

func deserializeValue(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case int:
	case float64:
	case string:
		return value, nil
	case *int:
	case *float64:
	case *string:
		return value, nil
	case jsonserialization.JsonParseNode:
	case *jsonserialization.JsonParseNode:
		return deserializeNode(v)
	case []*jsonserialization.JsonParseNode:
		return deserializeNodes(v)
	case []jsonserialization.JsonParseNode:
		return deserializeNodes(abs.CollectionApply(v, func(x jsonserialization.JsonParseNode) *jsonserialization.JsonParseNode {
			return &x
		}))
	default:
		return value, nil
	}
	return nil, nil
}

func castMapOfStrings(rawVal interface{}) (map[string]string, error) {
	result := make(map[string]string)
	v := reflect.ValueOf(rawVal)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			val, err := deserializeValue(v.MapIndex(key).Interface())
			if err != nil {
				return nil, err
			}
			result[key.String()] = *(val.(*string))
		}
	}
	return result, nil
}

func CreateBatchRequestItemDiscriminator(serialization.ParseNode) (serialization.Parsable, error) {
	var res batchItem
	return &res, nil
}
