package fileupload

import (
	"context"
	"fmt"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphgocore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

type UploadRequest[T interface{}] struct {
	requestAdapter     *msgraphgocore.GraphRequestAdapterBase
	data               []byte
	sessionUrlTemplate string
	rangeBegin         uint64
	rangeEnd           uint64
	rangeLength        uint64
	totalSessionLength uint64
	errorMappings      abstractions.ErrorMappings
}

func NewUploadRequest[T interface{}](requestAdapter *msgraphgocore.GraphRequestAdapterBase, data []byte, sessionUrlTemplate string, rangeBegin uint64, rangeEnd uint64, totalSessionLength uint64) *UploadRequest[T] {
	return &UploadRequest[T]{
		requestAdapter:     requestAdapter,
		data:               data,
		sessionUrlTemplate: sessionUrlTemplate,
		rangeBegin:         rangeBegin,
		rangeEnd:           rangeEnd,
		totalSessionLength: totalSessionLength,
	}
}

func (u *UploadRequest[T]) RangeLength() uint64 {
	return u.rangeEnd - u.rangeBegin + 1
}

func (u *UploadRequest[T]) UploadAsync() UploadResult[T] {
	uploadResult := UploadResult[T]{}
	requestInfo := u.createRequestInformation()
	err := u.requestAdapter.SendNoContent(context.Background(), requestInfo, u.errorMappings)

	return uploadResult
}

func (u *UploadRequest[T]) createRequestInformation() *abstractions.RequestInformation {
	requestInfo := abstractions.NewRequestInformation()
	requestInfo.UrlTemplate = u.sessionUrlTemplate
	requestInfo.Method = abstractions.POST
	requestInfo.SetStreamContent(u.data)

	requestInfo.Headers.Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", u.rangeLength, u.rangeEnd, u.totalSessionLength))
	requestInfo.Headers.Add("Content-Length", fmt.Sprintf("%d", u.RangeLength()))

	return requestInfo
}
