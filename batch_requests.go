package msgraphgocore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go-core/internal"
)

// Send sends a batch request
func Send(batch *batchRequest, adapter abstractions.RequestAdapter) (*BatchResponse, error) {
	batchJsonBody, err := batch.toJson()
	if err != nil {
		return nil, err
	}

	baseUrl, err := getBaseUrl(adapter)
	if err != nil {
		return nil, err
	}

	requestInfo := buildRequestInfo(batchJsonBody, baseUrl)
	return sendBatchRequest(requestInfo, adapter)
}

// AppendBatchItem converts RequestInformation to a BatchItem and adds it to a BatchRequest
//
// You can add upto 20 BatchItems to a BatchRequest
func (br *batchRequest) AppendBatchItem(reqInfo abstractions.RequestInformation) (*batchItem[serialization.Parsable], error) {
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
func (bi *batchItem[T]) DependsOnItem(item batchItem[serialization.Parsable]) {
	// DependsOn is a single value slice
	bi.DependsOn = []string{item.Id}
}

// GetBatchResponseById returns the response of the batch request item with the given id.
func GetBatchResponseById[T serialization.Parsable](resp *BatchResponse, itemId string) (T, error) {
	var res T

	for _, resp := range resp.Responses {
		if resp.Id == itemId {
			hasError := resp.Status >= 400 && resp.Status < 600

			if hasError {
				var errResp errorResponse

				jsonStr, err := json.Marshal(resp.Body)
				if err != nil {
					return res, err
				}
				err = json.Unmarshal(jsonStr, &errResp)
				if err != nil {
					return res, err
				}
				return res, errResp.Error
			}

			jsonStr, err := json.Marshal(resp.Body)
			if err != nil {
				return res, err
			}
			err = json.Unmarshal(jsonStr, &res)
			if err != nil {
				return res, err
			}

			return res, nil
		}
	}

	return res, errors.New("Response not found, check if id is valid")
}

// NewBatchRequest creates an instance of BatchRequest
func NewBatchRequest() *batchRequest {
	return &batchRequest{}
}

func toBatchItem(requestInfo abstractions.RequestInformation) (*batchItem[serialization.Parsable], error) {
	uri, err := requestInfo.GetUri()
	if err != nil {
		return nil, err
	}

	var body anyResponseBody
	err = json.Unmarshal(requestInfo.Content, &body)
	if err != nil {
		return nil, err
	}

	return &batchItem[serialization.Parsable]{
		Id:        uuid.NewString(),
		Method:    requestInfo.Method.String(),
		Body:      body,
		Headers:   requestInfo.Headers,
		Url:       uri.Path,
		DependsOn: make([]string, 0),
	}, nil
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

func sendBatchRequest(requestInfo *abstractions.RequestInformation, adapter abstractions.RequestAdapter) (*BatchResponse, error) {
	if requestInfo == nil {
		return nil, errors.New("requestInfo cannot be nil")
	}

	// SendAsync type asserts HandlerFunc's return value into a parsable. We bypass SendAsync deserialization by returning a noop
	// parsable struct and directly marshalling the response into a struct.
	var res BatchResponse
	_, err := adapter.SendAsync(requestInfo, nil, func(response any, errorMappings abstractions.ErrorMappings) (any, error) {
		resp, ok := response.(*http.Response)
		if !ok {
			return nil, errors.New("Response type assertion failed")
		}

		if status := resp.StatusCode; status >= 400 && status < 600 {
			return nil, fmt.Errorf("Request failed with status: %d", status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &res)
		if err != nil {
			return nil, err
		}

		return internal.NewNoOpParsable(), nil
	}, nil)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

type batchRequest struct {
	Requests []*batchItem[serialization.Parsable] `json:"requests"`
}

type BatchResponse struct {
	Responses []batchItem[serialization.Parsable]
}

type errorDetails struct {
	Code    string
	Message string
}

func (e errorDetails) Error() string {
	return fmt.Sprintf("Code: %s \n Message: %s", e.Code, e.Message)
}

type errorResponse struct {
	Error errorDetails
}

type batchItem[T serialization.Parsable] struct {
	Id        string            `json:"id"`
	Method    string            `json:"method"`
	Url       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      T                 `json:"body"`
	DependsOn []string          `json:"dependsOn"`
	Status    int               `json:"status,omitempty"`
}

type anyResponseBody map[string]any

func (h anyResponseBody) Serialize(writer serialization.SerializationWriter) error {
	panic("Not supported")
}

func (h anyResponseBody) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	panic("Not supported")
}
