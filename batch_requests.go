package msgraphgocore

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	abstractions "github.com/microsoft/kiota/abstractions/go"
	"github.com/microsoft/kiota/abstractions/go/serialization"
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

	var b map[string]interface{}
	json.Unmarshal(requestInfo.Content, &b)

	return &batchItem{
		Id:        uuid.NewString(),
		Method:    requestInfo.Method.String(),
		Body:      b,
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

func SendBatch(adapter abstractions.RequestAdapter, batch batchRequest) (serialization.Parsable, error) {
	var result serialization.Parsable

	jsonBody, err := batch.toJson()
	if err != nil {
		return result, err
	}

	ur := "https://graph.microsoft.com/v1.0/$batch"
	uri, err := url.Parse(ur)

	requestInfo := abstractions.NewRequestInformation()
	requestInfo.SetStreamContent(jsonBody)
	requestInfo.Method = abstractions.POST
	requestInfo.SetUri(*uri)
	requestInfo.Headers = map[string]string{
		"Content-Type": "application/json",
	}

	result, err = adapter.SendAsync(requestInfo, BatchResponseFactory, nil, nil)
	if err != nil {
		return result, err
	}

	resp := result.(*BatchResponse)
	for i := 0; i < len(resp.Responses); i++ {
		fmt.Println(*resp.Responses[i].GetStatus(), *resp.Responses[i].GetId())
	}

	return result, nil
}
