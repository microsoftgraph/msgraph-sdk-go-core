package msgraphgocore

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/google/uuid"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

type batchItem struct {
	Id        string                 `json:"id"`
	Method    string                 `json:"method"`
	Url       string                 `json:"url"`
	Headers   map[string]string      `json:"headers"`
	Body      map[string]interface{} `json:"body"`
	DependsOn []string               `json:"dependsOn"`
}

func newBatchItem(requestInfo abstractions.RequestInformation) (*batchItem, error) {
	url, err := requestInfo.GetUri()
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
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
	// DependsOn is a single value slice.
	b.DependsOn = []string{item.Id}
}

type batchRequest struct {
	Requests []*batchItem `json:"requests"`
}

func NewBatchRequest() *batchRequest {
	return &batchRequest{}
}

func (r *batchRequest) AppendItem(req abstractions.RequestInformation) (*batchItem, error) {
	if len(r.Requests) > 20 {
		return nil, errors.New("batch items limit exceeded. BatchRequest has a limit of 20 batch items")
	}

	batchItem, err := newBatchItem(req)
	if err != nil {
		return nil, err
	}

	r.Requests = append(r.Requests, batchItem)
	return batchItem, nil
}

func (r batchRequest) toJson() ([]byte, error) {
	return json.Marshal(r)
}

func (r *BatchResponse) GetResponses() []BatchItem {
	return r.Responses
}

type BResponse interface {
	GetResponses() map[string]interface{}
}

func BatchResponseFactory(parseNode serialization.ParseNode) (serialization.Parsable, error) {
	return NewBatchResponse(), nil
}

func SendBatch(adapter abstractions.RequestAdapter, batch batchRequest) (*BatchResponse, error) {
	var result *BatchResponse

	jsonBody, err := batch.toJson()
	if err != nil {
		return result, err
	}

	baseUrl, err := url.Parse(adapter.GetBaseUrl())
	if err != nil {
		return nil, err
	}

	requestInfo := abstractions.NewRequestInformation()
	requestInfo.SetStreamContent(jsonBody)
	requestInfo.Method = abstractions.POST
	requestInfo.SetUri(*baseUrl)
	requestInfo.Headers = map[string]string{
		"Content-Type": "application/json",
	}

	res, err := adapter.SendAsync(requestInfo, BatchResponseFactory, nil, nil)
	if err != nil {
		return result, err
	}

	return res.(*BatchResponse), nil
}
