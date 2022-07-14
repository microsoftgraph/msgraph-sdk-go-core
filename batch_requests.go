package msgraphgocore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go-core/internal"
)

// SendBatch sends a batch request
func SendBatch(batch batchRequest, adapter abstractions.RequestAdapter) (BatchResponse, error) {
	var res BatchResponse

	batchJsonBody, err := batch.toJson()
	if weGotAnError(err) {
		return res, err
	}

	baseUrl, err := getBaseUrl(adapter)
	if weGotAnError(err) {
		return res, err
	}

	requestInfo := buildRequestInfo(batchJsonBody, baseUrl)
	return sendBatchRequest(requestInfo, adapter)
}

// AppendBatchItem converts RequestInformation to a BatchItem and adds it to a BatchRequest
//
// You can add upto 20 BatchItems to a BatchRequest
func (br *batchRequest) AppendBatchItem(reqInfo abstractions.RequestInformation) (*batchItem, error) {
	if len(br.Requests) > 19 {
		return nil, errors.New("Batch items limit exceeded. BatchRequest has a limit of 20 batch items")
	}

	batchItem, err := toBatchItem(reqInfo)
	if weGotAnError(err) {
		return nil, err
	}

	br.Requests = append(br.Requests, batchItem)
	return batchItem, nil
}

// DependsOnItem creates a dependency chain between BatchItems.If A depends on B, then B will be sent before B
// A batchItem can only depend on one other batchItem
// see: https://docs.microsoft.com/en-us/graph/known-issues#request-dependencies-are-limited
func (bi *batchItem) DependsOnItem(item batchItem) {
	// DependsOn is a single value slice
	bi.DependsOn = []string{item.Id}
}

// NewBatchRequest creates an instance of BatchRequest
func NewBatchRequest() *batchRequest {
	return &batchRequest{}
}

func toBatchItem(requestInfo abstractions.RequestInformation) (*batchItem, error) {
	url, err := requestInfo.GetUri()
	if weGotAnError(err) {
		return nil, err
	}

	var body map[string]any
	json.Unmarshal(requestInfo.Content, &body)

	return &batchItem{
		Id:        uuid.NewString(),
		Method:    requestInfo.Method.String(),
		Body:      body,
		Headers:   requestInfo.Headers,
		Url:       url.Path,
		DependsOn: make([]string, 0),
	}, nil
}

func (r batchRequest) toJson() ([]byte, error) {
	return json.Marshal(r)
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

func sendBatchRequest(requestInfo *abstractions.RequestInformation, adapter abstractions.RequestAdapter) (BatchResponse, error) {
	var res BatchResponse

	// we don't care about SendAsync's return value because we don't use it and it is a workaround
	// on SendAsync's insistence to have a parsable node as a return value
	_, err := adapter.SendAsync(requestInfo, nil, func(response any, errorMappings abstractions.ErrorMappings) (any, error) {
		resp, ok := response.(*http.Response)
		if !ok {
			return nil, errors.New("Response type assertion failed")
		}

		if status := resp.StatusCode; status >= 400 && status <= 600 {
			return res, fmt.Errorf("Request failed with status: %d", status)
		}

		body, err := io.ReadAll(resp.Body)
		if weGotAnError(err) {
			return res, err
		}

		json.Unmarshal(body, &res)

		// returning a Noop Parsable here because SendAsync type asserts HandlerFunc's return value
		return internal.NewNoOpParsable(), nil
	}, nil)

	if weGotAnError(err) {
		return res, err
	}

	return res, nil
}

type batchRequest struct {
	Requests []*batchItem `json:"requests"`
}

type BatchResponse struct {
	Responses []batchItem
}

type batchItem struct {
	Id        string            `json:"id"`
	Method    string            `json:"method"`
	Url       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      map[string]any    `json:"body"`
	DependsOn []string          `json:"dependsOn"`
}
