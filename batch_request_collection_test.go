package msgraphgocore

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewBatchRequestCollectionNoLimit(t *testing.T) {
	batch := NewBatchRequestCollection(reqAdapter)
	reqInfo := getRequestInfo()

	for i := 0; i < 20; i++ {
		_, err := batch.AddBatchRequestStep(*reqInfo)
		if err != nil {
			return
		}
	}

	_, err := batch.AddBatchRequestStep(*reqInfo)
	assert.Nil(t, err)
}

func TestBatchRequestCollectionReturnsBatchResponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := getDummyJSON()
		w.WriteHeader(200)
		fmt.Fprint(w, jsonResponse)
	}))
	defer testServer.Close()

	reqInfo := getRequestInfo()

	mockPath := testServer.URL + "/$batch"
	reqAdapter.SetBaseUrl(mockPath) // check that path is not empty instead

	batch := NewBatchRequestCollection(reqAdapter)
	for i := 0; i < 40; i++ {
		_, err := batch.AddBatchRequestStep(*reqInfo)
		if err != nil {
			require.NoError(t, err)
		}
	}

	resp, err := batch.Send(context.Background(), reqAdapter)
	require.NoError(t, err)

	assert.Equal(t, len(resp.GetResponses()), 12)
}

func TestBatchRequestResponseGetFailedResponses(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := getDummyJSON()
		w.WriteHeader(200)
		fmt.Fprint(w, jsonResponse)
	}))
	defer testServer.Close()

	reqInfo := getRequestInfo()

	mockPath := testServer.URL + "/$batch"
	reqAdapter.SetBaseUrl(mockPath) // check that path is not empty instead

	batch := NewBatchRequestCollection(reqAdapter)
	_, err := batch.AddBatchRequestStep(*reqInfo)
	require.NoError(t, err)

	resp, err := batch.Send(context.Background(), reqAdapter)
	require.NoError(t, err)

	assert.Equal(t, len(resp.GetStatusCodes()), 4)

	status := resp.GetFailedResponses()
	assert.Equal(t, 1, len(status))
}
