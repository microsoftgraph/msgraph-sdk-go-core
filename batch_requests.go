package msgraphgocore

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go-core/internal"
)

func SendBatch(batch batchRequest, adapter abstractions.RequestAdapter) (BatchResponse, error) {
	var res BatchResponse

	jsonBody, err := batch.toJson()
	if weGotAnError(err) {
		return res, err
	}

	baseUrl, err := getBaseUrl(adapter)
	if weGotAnError(err) {
		return res, err
	}

	requestInfo := buildRequestInfo(jsonBody, baseUrl)
	return sendBatchRequest(requestInfo, adapter)
}

func NewBatchRequest() *batchRequest {
	return &batchRequest{}
}

func (r *batchRequest) AppendBatchItem(reqInfo abstractions.RequestInformation) (*batchItem, error) {
	if len(r.Requests) > 20 {
		return nil, errors.New("batch items limit exceeded. BatchRequest has a limit of 20 batch items")
	}

	batchItem, err := toBatchItem(reqInfo)
	if weGotAnError(err) {
		return nil, err
	}

	r.Requests = append(r.Requests, batchItem)
	return batchItem, nil
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

func (b *batchItem) DependsOnItem(item batchItem) {
	// DependsOn is a single value slice
	b.DependsOn = []string{item.Id}
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
