package fileupload

import (
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go-core"
	"time"
)

type LargeFileUploadTask[T interface{}] struct {
	uploadSession  UploadSession
	requestAdapter *msgraphgocore.GraphRequestAdapterBase
	data           []byte
	sliceSize      int
}

type UploadSession interface {
	GetNextExpectedRanges() []string
	GetUploadUrl() *string
	GetExpirationDateTime() *time.Time
}

type ProgressCallback func(current float64, total float64)

type UploadResult[T interface{}] struct {
	ItemResponse    T
	UploadSession   UploadSession
	URI             string
	UploadSucceeded bool
	ErrorMappings   abstractions.ErrorMappings
}

const DefaultSliceSize = 1024

func NewLargeFileUploadTask[T interface{}](requestAdapter *msgraphgocore.GraphRequestAdapterBase, uploadSession UploadSession, data []byte, maxSliceSize int) *LargeFileUploadTask[T] {
	return &LargeFileUploadTask[T]{
		uploadSession:  uploadSession,
		requestAdapter: requestAdapter,
		data:           data,
		sliceSize:      maxSliceSize,
	}
}

func (l *LargeFileUploadTask[T]) Upload(progress ProgressCallback) UploadResult[T] {
	chunkSize := l.sliceSize
	if chunkSize == -1 {
		chunkSize = DefaultSliceSize
	}

	uploadResult := UploadResult[T]{}

	chunks := msgraphgocore.ChunkSlice(l.data, chunkSize)
	for _, dataSection := range chunks {
		uploadRequest := NewUploadRequest[T](l.requestAdapter, dataSection, l.uploadSession.GetUploadUrl())
		uploadRequest.UploadAsync()
	}
	return uploadResult
}
