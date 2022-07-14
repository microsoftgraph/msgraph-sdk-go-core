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

func SendBatch(adapter abstractions.RequestAdapter, batch batchRequest) (BatchResponse, error) {
	var res BatchResponse

	jsonBody, err := batch.toJson()
	if weGotAnError(err) {
		return res, err
	}

	baseUrl, err := url.Parse(adapter.GetBaseUrl())
	if weGotAnError(err) {
		return res, err
	}

	requestInfo := abstractions.NewRequestInformation()
	requestInfo.SetStreamContent(jsonBody)
	requestInfo.Method = abstractions.POST
	requestInfo.SetUri(*baseUrl)
	requestInfo.Headers = map[string]string{
		"Content-Type": "application/json",
	}

	adapter.SendAsync(requestInfo, nil, func(response interface{}, errorMappings abstractions.ErrorMappings) (interface{}, error) {
		resp, _ := response.(*http.Response)
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &res)

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

func NewBatchRequest() *batchRequest {
	return &batchRequest{}
}

func newBatchItem(requestInfo abstractions.RequestInformation) (*batchItem, error) {
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

func (r *batchRequest) AppendBatchItem(req abstractions.RequestInformation) (*batchItem, error) {
	if len(r.Requests) > 20 {
		return nil, errors.New("batch items limit exceeded. BatchRequest has a limit of 20 batch items")
	}

	batchItem, err := newBatchItem(req)
	if weGotAnError(err) {
		return nil, err
	}

	r.Requests = append(r.Requests, batchItem)
	return batchItem, nil
}

func (r batchRequest) toJson() ([]byte, error) {
	return json.Marshal(r)
}
