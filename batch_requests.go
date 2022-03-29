package msgraphgocore

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	abstractions "github.com/microsoft/kiota/abstractions/go"
)

type batchItem struct {
	Id        string            `json:"id"`
	Method    string            `json:"method"`
	Url       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	DependsOn []string          `json:"dependsOn"`
}

func newBatchItem(requestInfo abstractions.RequestInformation) (*batchItem, error) {
	url, err := requestInfo.GetUri()
	if err != nil {
		return nil, err
	}

	return &batchItem{
		Id:        uuid.NewString(),
		Method:    requestInfo.Method.String(),
		Body:      string(requestInfo.Content),
		Headers:   requestInfo.Headers,
		Url:       url.Path,
		DependsOn: []string{},
	}, nil
}

func (b *batchItem) dependsOn(item batchItem) {
	// DependsOn is a single value slice.
	b.DependsOn[0] = item.Id
}

type BatchRequest struct {
	Requests []*batchItem `json:"requests"`
}

func NewBatchRequest() *BatchRequest {
	return &BatchRequest{}
}

func (r *BatchRequest) appendItem(req abstractions.RequestInformation) (*batchItem, error) {
	if len(r.Requests) > 20 {
		return nil, errors.New("Batch items limit exceeded. BatchRequest has a limit of 20 batch items")
	}

	batchItem, err := newBatchItem(req)
	if err != nil {
		return nil, err
	}

	r.Requests = append(r.Requests, batchItem)
	return batchItem, nil
}

func (r BatchRequest) toJson() ([]byte, error) {
	return json.Marshal(r)
}

func SendBatch(client http.Client, batch BatchRequest) (string, error) {
	var result string

	jsonBody, err := batch.toJson()
	if err != nil {
		return result, err
	}

	res, err := client.Post(
		"https://graph.microsoft.com/v1.0/$batch",
		"application/json",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return result, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	json.Unmarshal(b, &result)
	return result, nil
}
