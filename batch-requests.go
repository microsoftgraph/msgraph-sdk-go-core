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

type BatchItem struct {
	Id        string
	Method    abstractions.HttpMethod
	Body      []byte
	Headers   map[string]string
	Url       string
	DependsOn []string
}

func newBatchItem(requestInfo abstractions.RequestInformation) (*BatchItem, error) {
	url, err := requestInfo.GetUri()
	if err != nil {
		return nil, err
	}

	return &BatchItem{
		Id:      uuid.NewString(),
		Method:  requestInfo.Method,
		Body:    requestInfo.Content,
		Headers: requestInfo.Headers,
		Url:     url.String(),
	}, nil
}

func (b *BatchItem) dependsOn(item BatchItem) {
	// DependsOn is a single value slice.
	b.DependsOn[0] = item.Id
}

type BatchRequest struct {
	Requests []BatchItem
}

func NewBatchRequest(requestItems []BatchItem) *BatchRequest {
	return &BatchRequest{
		Requests: requestItems,
	}
}

func (r *BatchRequest) appendItem(req abstractions.RequestInformation) (*BatchItem, error) {
	if len(r.Requests) > 20 {
		return nil, errors.New("Batch items limit exceeded. BatchRequest has a limit of 20 batch items")
	}

	batchItem, err := newBatchItem(req)
	if err != nil {
		return nil, err
	}

	r.Requests = append(r.Requests, *batchItem)
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
