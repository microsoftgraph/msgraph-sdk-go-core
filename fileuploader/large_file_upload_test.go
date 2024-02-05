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

	uploader := NewLargeFileUploadTask[internal.UploadResponseble](reqAdapter, uploadSession, byteStream, int64(maxSliceSize), internal.CreateUploadResponseFromDiscriminatorValue, errorMapping)

	// verify that the object was created correctly
	// verify the number of sub upload tasks

	progressCall := 0
	progress := func(progress int64, total int64) {
		progressCall++
	}
	result := uploader.UploadAsync(progress)

	// verify that status is correct
	assert.True(t, result.GetUploadSucceeded())
	assert.Equal(t, 12, progressCall) // progress callback should be called for every sub upload task
}

type mockUploadSession struct {
	UploadUrl          string
	ExpectedRanges     []string
	OdataType          string
	ExpirationDateTime time.Time
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
