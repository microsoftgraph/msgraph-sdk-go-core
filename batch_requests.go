package msgraphgocore

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

const BATCH_REQUEST_ERROR_REGISTRY_KEY = "BATCH_REQUEST_ERROR_REGISTRY_KEY"

// Send sends a batch request
func (br *batchRequest) Send(ctx context.Context, adapter abstractions.RequestAdapter) (BatchResponse, error) {
	batchJsonBody, err := br.toJson()
	if err != nil {
		return nil, err
	}

	baseUrl, err := getBaseUrl(adapter)
	if err != nil {
		return nil, err
	}

	requestInfo := buildRequestInfo(batchJsonBody, baseUrl)
	return sendBatchRequest(ctx, requestInfo, adapter)
}

func (br *batchRequest) toJson() ([]byte, error) {
	return json.Marshal(br)
}

func getBaseUrl(adapter abstractions.RequestAdapter) (*url.URL, error) {
	return url.Parse(adapter.GetBaseUrl())
}

func buildRequestInfo(jsonBody []byte, baseUrl *url.URL) *abstractions.RequestInformation {
	requestInfo := abstractions.NewRequestInformation()

	requestInfo.SetStreamContent(jsonBody)
	requestInfo.Method = abstractions.POST
	requestInfo.SetUri(*baseUrl)
	requestInfo.Headers = map[string]string{
		"Content-Type": "application/json",
	}

	return requestInfo
}

// AppendBatchItem converts RequestInformation to a BatchItem and adds it to a BatchRequest
//
// You can add upto 20 BatchItems to a BatchRequest
func (br *batchRequest) AppendBatchItem(reqInfo abstractions.RequestInformation) (BatchItem, error) {
	if len(br.Requests) > 19 {
		return nil, errors.New("Batch items limit exceeded. BatchRequest has a limit of 20 batch items")
	}

	batchItem, err := toBatchItem(reqInfo)
	if err != nil {
		return nil, err
	}

	br.Requests = append(br.Requests, batchItem)
	return batchItem, nil
}

// DependsOnItem creates a dependency chain between BatchItems.If A depends on B, then B will be sent before B
// A batchItem can only depend on one other batchItem
// see: https://docs.microsoft.com/en-us/graph/known-issues#request-dependencies-are-limited
func (bi *batchItem) DependsOnItem(item BatchItem) {
	// DependsOn is a errorRegistry value slice
	bi.DependsOn = []string{*item.GetId()}
}

func getResponsePrimaryContentType(responseItem BatchItem) string {
	header := responseItem.GetHeaders()
	if header == nil {
		return ""
	}
	rawType := header["Content-Type"]
	splat := strings.Split(rawType, ";")
	return strings.ToLower(splat[0])
}

func getRootParseNode(responseItem BatchItem) (absser.ParseNode, error) {
	contentType := getResponsePrimaryContentType(responseItem)
	if contentType == "" {
		return nil, nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(responseItem.GetBody())
	if err != nil {
		return nil, err
	}
	return serialization.DefaultParseNodeFactoryInstance.GetRootParseNode(contentType, buf.Bytes())
}

func throwErrors(responseItem BatchItem, typeName string) error {
	errorMappings := getErrorMapper(typeName)
	responseStatus := *responseItem.GetStatus()

	statusAsString := strconv.Itoa(int(responseStatus))
	var errorCtor absser.ParsableFactory = nil
	if len(errorMappings) != 0 {
		if responseStatus >= 400 && responseStatus < 500 && errorMappings["4XX"] != nil {
			errorCtor = errorMappings["4XX"]
		} else if responseStatus >= 500 && responseStatus < 600 && errorMappings["5XX"] != nil {
			errorCtor = errorMappings["5XX"]
		}
	}

	if errorCtor == nil {
		return &abstractions.ApiError{
			Message: "The server returned an unexpected status code and no error factory is registered for this code: " + statusAsString,
		}
	}

	rootNode, err := getRootParseNode(responseItem)
	if err != nil {
		return err
	}
	if rootNode == nil {
		return &abstractions.ApiError{
			Message: "The server returned an unexpected status code with no response body: " + statusAsString,
		}
	}

	errValue, err := rootNode.GetObjectValue(errorCtor)
	if err != nil {
		return err
	}

	return errValue.(error)
}

// GetBatchResponseById returns the response of the batch request item with the given id.
func GetBatchResponseById[T serialization.Parsable](resp BatchResponse, itemId string) (*T, error) {
	var res T
	item := resp.GetResponseById(itemId)

	hasError := *item.GetStatus() >= 400 && *item.GetStatus() < 600
	if hasError {
		typeName := reflect.TypeOf(res).Name()
		return nil, throwErrors(item, typeName)
	}

	jsonStr, err := json.Marshal(item.GetBody())
	if err != nil {
		return &res, err
	}
	err = json.Unmarshal(jsonStr, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

// NewBatchRequest creates an instance of BatchRequest
func NewBatchRequest() *batchRequest {
	return &batchRequest{}
}

func toBatchItem(requestInfo abstractions.RequestInformation) (*batchItem, error) {
	uri, err := requestInfo.GetUri()
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	err = json.Unmarshal(requestInfo.Content, &body)
	if err != nil {
		return nil, err
	}

	newID := uuid.NewString()
	method := requestInfo.Method.String()

	return &batchItem{
		Id:        &newID,
		method:    &method,
		Body:      body,
		Headers:   requestInfo.Headers,
		Url:       &uri.Path,
		DependsOn: make([]string, 0),
	}, nil
}

func getErrorMapper(key string) abstractions.ErrorMappings {
	errorMapperSrc, found := GetErrorFactoryFromRegistry(key)
	if found {
		return errorMapperSrc
	}
	return nil
}

func sendBatchRequest(ctx context.Context, requestInfo *abstractions.RequestInformation, adapter abstractions.RequestAdapter) (BatchResponse, error) {
	if requestInfo == nil {
		return nil, errors.New("requestInfo cannot be nil")
	}

	response, err := adapter.SendAsync(ctx, requestInfo, CreateBatchResponseDiscriminator, getErrorMapper(BATCH_REQUEST_ERROR_REGISTRY_KEY))
	if err != nil {
		return nil, err
	}

	return response.(BatchResponse), nil
}
