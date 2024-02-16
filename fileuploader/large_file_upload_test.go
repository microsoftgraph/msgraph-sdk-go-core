package fileuploader

import (
	"fmt"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	msgraphgocore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go-core/internal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func prepareUploader(testServer *httptest.Server) LargeFileUploadTask[internal.UploadResponseble] {
	absser.DefaultParseNodeFactoryInstance.ContentTypeAssociatedFactories["application/json"] = jsonserialization.NewJsonParseNodeFactory()

	reqAdapter, _ := msgraphgocore.NewGraphRequestAdapterBase(&authentication.AnonymousAuthenticationProvider{}, msgraphgocore.GraphClientOptions{
		GraphServiceVersion:        "",
		GraphServiceLibraryVersion: "",
	})

	mockPath := testServer.URL + "/uploadUrl"
	reqAdapter.SetBaseUrl(mockPath)

	byteStream := &internal.MockByteStream{
		Content: []byte("mock byteStream content"),
	}

	uploadSession := &mockUploadSession{
		UploadUrl:          mockPath,
		ExpectedRanges:     []string{"0-4", "6-"},
		OdataType:          "odatatype",
		ExpirationDateTime: time.Time{},
	}
	maxSliceSize := 2

	errorMapping := abstractions.ErrorMappings{
		"4XX": internal.CreateSampleErrorFromDiscriminatorValue,
		"5XX": internal.CreateSampleErrorFromDiscriminatorValue,
	}

	return NewLargeFileUploadTask[internal.UploadResponseble](reqAdapter, uploadSession, byteStream, int64(maxSliceSize), internal.CreateUploadResponseFromDiscriminatorValue, errorMapping)
}

func TestLargeFileUploadTask(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := `{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#microsoft.graph.uploadSession",
			"uploadUrl": "https://uploadUrl",
			"expirationDateTime": "2021-08-10T00:00:00Z"
		}`
		w.WriteHeader(200)
		fmt.Fprint(w, jsonResponse)
	}))
	defer testServer.Close()

	uploader := prepareUploader(testServer)

	// verify that the object was created correctly
	// verify the number of sub upload tasks
	progressCall := 0
	progress := func(progress int64, total int64) {
		progressCall++
	}
	result := uploader.Upload(progress)

	// verify that status is correct
	assert.True(t, result.GetUploadSucceeded())
	assert.Equal(t, 12, progressCall) // progress callback should be called for every sub upload task
}

func TestResumeLargeFileUploadTask(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		testTime := time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z")
		jsonResponse := `{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#microsoft.graph.uploadSession",
			"uploadUrl": "https://uploadUrl",
			"expirationDateTime": "%s",
			"nextExpectedRanges": ["0-4", "6-"]
		}`
		w.WriteHeader(200)
		formattedResponse := fmt.Sprintf(jsonResponse, testTime)
		fmt.Fprint(w, formattedResponse)
	}))
	defer testServer.Close()

	uploader := prepareUploader(testServer)

	progressCall := 0
	progress := func(progress int64, total int64) {
		progressCall++
	}
	result, err := uploader.Resume(progress)
	assert.NoError(t, err)

	// verify that status is correct
	assert.True(t, result.GetUploadSucceeded())
	assert.Equal(t, 12, progressCall) // progress callback should be called for every sub upload task

}

func TestCancelLargeFileUploadTask(t *testing.T) {

	var receivedReq *http.Request
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(204)
		receivedReq = req
	}))
	defer testServer.Close()

	uploader := prepareUploader(testServer)
	err := uploader.Cancel()
	assert.NoError(t, err)
	assert.Equal(t, "DELETE", receivedReq.Method)
}

type mockUploadSession struct {
	UploadUrl          string
	ExpectedRanges     []string
	OdataType          string
	ExpirationDateTime time.Time
}

func (m *mockUploadSession) SetExpirationDateTime(expirationDateTime *time.Time) {
	m.ExpirationDateTime = *expirationDateTime
}

func (m *mockUploadSession) SetNextExpectedRanges(nextExpectedRanges []string) {
	m.ExpectedRanges = nextExpectedRanges
}

func (m *mockUploadSession) Serialize(writer absser.SerializationWriter) error {
	return nil
}

func (m *mockUploadSession) GetFieldDeserializers() map[string]func(absser.ParseNode) error {
	return make(map[string]func(absser.ParseNode) error)
}

func (m *mockUploadSession) GetExpirationDateTime() *time.Time {
	return &m.ExpirationDateTime
}

func (m *mockUploadSession) GetNextExpectedRanges() []string {
	return m.ExpectedRanges
}

func (m *mockUploadSession) GetOdataType() *string {
	return &m.OdataType
}

func (m *mockUploadSession) GetUploadUrl() *string {
	return &m.UploadUrl
}
