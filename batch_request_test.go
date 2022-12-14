package msgraphgocore

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-sdk-go-core/internal"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func p[T interface{}](t T) *T {
	return &t
}

func TestConstructionOfRequests(t *testing.T) {
	reqInfo := getRequestInfo()

	batch := NewBatchRequest()

	item1, err := batch.AddBatchRequestStep(*reqInfo)
	require.NoError(t, err)

	item2, err := batch.AddBatchRequestStep(*reqInfo)
	require.NoError(t, err)

	assert.Equal(t, len(batch.GetRequests()), 2)
	assert.Equal(t, batch.GetRequests()[0], item1)
	assert.Equal(t, batch.GetRequests()[1], item2)
}

func TestRegisteringDependsOn(t *testing.T) {

	reqInfo1 := getRequestInfo()
	reqInfo2 := getRequestInfo()

	batch := NewBatchRequest()
	batchItem1, err := batch.AddBatchRequestStep(*reqInfo1)
	require.NoError(t, err)

	batchItem2, err := batch.AddBatchRequestStep(*reqInfo2)
	require.NoError(t, err)

	batchItem2.DependsOnItem(batchItem1)

	assert.Equal(t, batchItem2.GetDependsOn(), []string{*batchItem1.GetId()})
}

func TestReturnsBatchResponse(t *testing.T) {
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

	batch := NewBatchRequest()
	_, err := batch.AddBatchRequestStep(*reqInfo)
	require.NoError(t, err)

	resp, err := batch.Send(context.Background(), reqAdapter)
	require.NoError(t, err)

	assert.Equal(t, len(resp.GetResponses()), 4)
}

func TestContentSentToServer(t *testing.T) {
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

	batch := NewBatchRequest()
	item, err := batch.AddBatchRequestStep(*reqInfo)
	item.SetId(p("123"))
	require.NoError(t, err)

	baseUrl, err := getBaseUrl(reqAdapter)
	require.NoError(t, err)

	requestInfo, err := buildRequestInfo(context.Background(), reqAdapter, batch, baseUrl)
	require.NoError(t, err)
	content := string(requestInfo.Content)
	expected := "{\"requests\":[{\"id\":\"123\",\"method\":\"GET\",\"url\":\"\",\"headers\":{\"content-type\":\"application/json\"},\"body\":{\"username\":\"name\"},\"dependsOn\":[]}]}"
	assert.Equal(t, expected, content)

	resp, err := batch.Send(context.Background(), reqAdapter)
	require.NoError(t, err)

	assert.Equal(t, len(resp.GetResponses()), 4)
}

func TestRespectsBatchItemLimitOf20BatchItems(t *testing.T) {
	batch := NewBatchRequest()
	reqInfo := getRequestInfo()

	for i := 0; i < 20; i++ {
		_, err := batch.AddBatchRequestStep(*reqInfo)
		if err != nil {
			return
		}
	}

	_, err := batch.AddBatchRequestStep(*reqInfo)
	assert.Equal(t, err.Error(), "batch items limit exceeded. BatchRequest has a limit of 20 batch items")
}

func TestHandlesUnhandledHTTPError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(403)
		fmt.Fprint(w, "")
	}))
	defer testServer.Close()

	mockPath := testServer.URL + "/$batch"
	reqAdapter.SetBaseUrl(mockPath)

	reqInfo := getRequestInfo()
	batch := NewBatchRequest()
	_, err := batch.AddBatchRequestStep(*reqInfo)
	require.NoError(t, err)

	_, err = batch.Send(context.Background(), reqAdapter)
	assert.Equal(t, err.Error(), "The server returned an unexpected status code and no error factory is registered for this code: 403")
}

func TestHandlesHTTPError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(403)
		fmt.Fprint(w, "")
	}))
	defer testServer.Close()

	mockPath := testServer.URL + "/$batch"
	reqAdapter.SetBaseUrl(mockPath)

	reqInfo := getRequestInfo()
	batch := NewBatchRequest()
	batch.AddBatchRequestStep(*reqInfo)

	errorMapping := abstractions.ErrorMappings{
		"4XX": internal.CreateSampleErrorFromDiscriminatorValue,
		"5XX": internal.CreateSampleErrorFromDiscriminatorValue,
	}
	// register errorMapper
	err := RegisterError(BatchRequestErrorRegistryKey, errorMapping)
	require.NoError(t, err)

	_, err = batch.Send(context.Background(), reqAdapter)
	assert.Equal(t, err.Error(), "content is empty")

	err = DeRegisterError(BatchRequestErrorRegistryKey)
	require.NoError(t, err)
}

func TestGetResponseByIdForSuccessfulRequest(t *testing.T) {
	mockResponse := `{
			"responses": [
			  {
				"id": "2",
				"status": 200,
				"body": {
				  "username": "testuser",
				  "person" : {
					"firstName" : "Tony",
					"lastName" : "Blair",
					"active" : false,
					"bankBalance": 234234.67,
					"accounts" : [1,2,3],
					"positions" : ["Prime","Minister"],
					"children" : [
					  {
						"firstName" : "Kathryn",
						"lastName" : "Blair"
					  },
					  {
						"firstName" : "Euan",
						"lastName" : "Blair"
					  }
					]
				  }
				}
			  }
			]
		  }`
	mockServer := makeMockRequest(200, mockResponse)
	defer mockServer.Close()

	mockPath := mockServer.URL + "/$batch"
	reqAdapter.SetBaseUrl(mockPath)

	reqInfo := getRequestInfo()
	batch := NewBatchRequest()
	_, err := batch.AddBatchRequestStep(*reqInfo)
	if err != nil {
		return
	}

	resp, err := batch.Send(context.Background(), reqAdapter)
	require.NoError(t, err)

	user, err := GetBatchResponseById[User](resp, "2")
	require.NoError(t, err)

	assert.Equal(t, user.UserName, "testuser")
}

type User struct {
	UserName string `json:"username"`
	Person   Person `json:"person"`
}

func (u User) Serialize(writer serialization.SerializationWriter) error {
	panic("implement me")
}

func (u User) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	panic("implement me")
}

type Person struct {
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Active      bool      `json:"active"`
	Positions   []*string `json:"positions"`
	BankBalance *float64  `json:"bankBalance"`
	Accounts    []*int    `json:"accounts"`
	Children    []*Person `json:"children"`
}

func (u Person) Serialize(writer serialization.SerializationWriter) error {
	panic("implement me")
}

func (u Person) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	panic("implement me")
}

func TestGetResponseByIdFailedRequest(t *testing.T) {
	mockServer := makeMockRequest(200, getDummyJSON())
	defer mockServer.Close()

	mockPath := mockServer.URL + "/$batch"
	reqAdapter.SetBaseUrl(mockPath)

	reqInfo := getRequestInfo()
	batch := NewBatchRequest()
	_, err := batch.AddBatchRequestStep(*reqInfo)
	require.NoError(t, err)

	resp, err := batch.Send(context.Background(), reqAdapter)
	require.NoError(t, err)

	_, err = GetBatchResponseById[User](resp, "3")
	assert.Equal(t, "The server returned an unexpected status code and no error factory is registered for this code: 401", err.Error())
}

func TestGetResponseByIdFailedRequestWithFactory(t *testing.T) {
	mockServer := makeMockRequest(200, getDummyJSON())
	defer mockServer.Close()

	mockPath := mockServer.URL + "/$batch"
	reqAdapter.SetBaseUrl(mockPath)

	errorMapping := abstractions.ErrorMappings{
		"4XX": internal.CreateSampleErrorFromDiscriminatorValue,
		"5XX": internal.CreateSampleErrorFromDiscriminatorValue,
	}
	// register errorMapper
	err := RegisterError(BatchRequestErrorRegistryKey, errorMapping)
	require.NoError(t, err)

	reqInfo := getRequestInfo()
	batch := NewBatchRequest()
	_, err = batch.AddBatchRequestStep(*reqInfo)
	require.NoError(t, err)

	resp, err := batch.Send(context.Background(), reqAdapter)
	require.NoError(t, err)

	_, err = GetBatchResponseById[User](resp, "3")
	assert.Equal(t, "The server returned an unexpected status code with no response body: 401", err.Error())

	err = DeRegisterError(BatchRequestErrorRegistryKey)
	require.NoError(t, err)
}

func makeMockRequest(mockStatus int, mockResponse string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(mockStatus)
		fmt.Fprint(w, mockResponse)
	}))
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
	reqInfo.UrlTemplate = "{+baseurl}/$batch"
	headers := abstractions.NewRequestHeaders()
	headers.Add("Content-Type", "application/json")
	reqInfo.Headers.AddAll(headers)

	return reqInfo
}

func getDummyJSON() string {
	return `{
	"responses": [
	{
	  "id": "1",
	  "status": 302,
	  "body": null,
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
	      "message": "Insufficient permissions"
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
	  "url": "https://graph.microsoft.com/v1.0/$metadata#Collection(microsoft.graph.plannerTask)",
	  "body": null
	}
	]
				}`
}
