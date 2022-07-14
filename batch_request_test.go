package msgraphgocore

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratesJSONFromRequestBody(t *testing.T) {
	reqInfo := getRequestInfo()

	batch := NewBatchRequest()
	item, _ := batch.AppendBatchItem(*reqInfo)
	item.Id = "1"

	expected := "{\"requests\":[{\"id\":\"1\",\"method\":\"GET\",\"url\":\"\",\"headers\":{\"content-type\":\"application/json\"},\"body\":{\"username\":\"name\"},\"dependsOn\":[]}]}"
	actual, _ := batch.toJson()

	assert.Equal(t, expected, string(actual))
}

func TestDependsOnRelationshipInBatchRequestItems(t *testing.T) {

	reqInfo1 := getRequestInfo()
	reqInfo2 := getRequestInfo()

	batch := NewBatchRequest()
	batchItem1, _ := batch.AppendBatchItem(*reqInfo1)
	batchItem2, _ := batch.AppendBatchItem(*reqInfo2)
	batchItem1.Id = "1"
	batchItem2.Id = "2"

	batchItem2.DependsOnItem(*batchItem1)

	expected := "{\"requests\":[{\"id\":\"1\",\"method\":\"GET\",\"url\":\"\",\"headers\":{\"content-type\":\"application/json\"},\"body\":{\"username\":\"name\"},\"dependsOn\":[]},{\"id\":\"2\",\"method\":\"GET\",\"url\":\"\",\"headers\":{\"content-type\":\"application/json\"},\"body\":{\"username\":\"name\"},\"dependsOn\":[\"1\"]}]}"
	actual, err := batch.toJson()
	require.NoError(t, err)

	assert.Equal(t, expected, string(actual))
}

func TestReturnsBatchResponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		jsonResponse := getDummyJSON()
		w.WriteHeader(200)
		fmt.Fprint(w, jsonResponse)
	}))
	defer testServer.Close()

	mockPath := testServer.URL + "/$batch"
	reqAdapter.SetBaseUrl(mockPath)

	reqInfo := getRequestInfo()
	batch := NewBatchRequest()
	batch.AppendBatchItem(*reqInfo)

	resp, err := SendBatch(*batch, reqAdapter)
	require.NoError(t, err)

	assert.Equal(t, len(resp.Responses), 4)
}

func getRequestInfo() *abstractions.RequestInformation {
	content := `
{
    "username": "name"
}
`
	reqInfo := abstractions.NewRequestInformation()
	reqInfo.SetUri(url.URL{})
	reqInfo.Content = []byte(content)
	reqInfo.Headers = map[string]string{"content-type": "application/json"}

	return reqInfo
}

func getDummyJSON() string {
	return `{
	"responses": [
	{
	  "id": "1",
	  "status": 302,
	  "headers": {
	    "location": "https://b0mpua-by3301.files.1drv.com/y23vmagahszhxzlcvhasdhasghasodfi"
	  }
	},
	{
	  "id": "3",
	  "status": 401,
	  "body": {
	    "error": {
	      "code": "Forbidden",
	      "message": "..."
	    }
	  }
	},
	{
	  "id": "2",
	  "status": 200,
	  "body": {
	    "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#Collection(microsoft.graph.plannerTask)",
	    "value": []
	  }
	},
	{
	  "id": "4",
	  "status": 204,
	  "body": null
	}
	]
				}`
}
