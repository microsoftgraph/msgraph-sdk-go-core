package fileuploader

import (
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-sdk-go-core/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLargeFileUploadTask(t *testing.T) {
	byteStream := &internal.MockByteStream{
		Content: []byte("mock byteStream content"),
	}

	var adapter abstractions.RequestAdapter
	var uploadSession UploadSession
	maxSliceSize := 8000

	errorMapping := abstractions.ErrorMappings{
		"4XX": internal.CreateSampleErrorFromDiscriminatorValue,
		"5XX": internal.CreateSampleErrorFromDiscriminatorValue,
	}

	uploader := NewLargeFileUploadTask[internal.User](adapter, uploadSession, byteStream, int64(maxSliceSize), CreateDriveItemFromDiscriminatorValue, errorMapping)

	// verify that the object was created correctly
	// verify the number of sub upload tasks

	progressCall := 0
	progress := func(progress int64, total int64) {
		progressCall++
	}
	result := uploader.UploadAsync(progress)

	// verify that status is correct
	assert.True(t, result.GetUploadSucceeded())
	assert.Equal(t, progressCall, 1)
}

func CreateDriveItemFromDiscriminatorValue(parseNode serialization.ParseNode) (serialization.Parsable, error) {
	res := internal.SampleError{}
	return &res, nil
}
