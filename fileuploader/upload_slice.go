package fileuploader

import (
	"context"
	"fmt"
	abstractions "github.com/microsoft/kiota-abstractions-go"
)

type uploadSlice[T interface{}] struct {
	RequestAdapter     abstractions.RequestAdapter
	UrlTemplate        string
	RangeBegin         float64
	RangeEnd           float64
	TotalSessionLength float64
	RangeLength        float64
	data               []byte
	errorMappings      abstractions.ErrorMappings
}

func (l *largeFileUploadTask[T]) createUploadSlices() []uploadSlice[T] {
	rangesRemaining := l.getRangesRemaining()

	uploadSlices := make([]uploadSlice[T], len(rangesRemaining))

	for i, v := range rangesRemaining {
		uploadSlices[i] = uploadSlice[T]{
			RequestAdapter: l.adapter,
			UrlTemplate:    *l.uploadSession.GetUploadUrl(),
			RangeBegin:     v.Start,
			RangeEnd:       v.End,
		}
	}

	return uploadSlices
}

func (u *uploadSlice[T]) UploadAsync() (UploadResult[T], error) {
	res := NewUploadResult[T]()
	requestInfo := u.createRequestInformation(u.data)

	var uploadResponseHandler abstractions.ResponseHandler = func(response interface{}, errorMappings abstractions.ErrorMappings) (interface{}, error) {
		panic("To do")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, abstractions.ResponseHandlerOptionKey, uploadResponseHandler)

	err := u.RequestAdapter.SendNoContent(ctx, requestInfo, u.errorMappings)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *uploadSlice[T]) createRequestInformation(content []byte) *abstractions.RequestInformation {
	headers := abstractions.NewRequestHeaders()
	headers.Add("Content-Range", fmt.Sprintf("bytes %f-%f/%f", u.RangeLength, u.RangeEnd, u.TotalSessionLength))
	headers.Add("Content-Length", fmt.Sprintf("%f", u.RangeLength))

	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Headers = headers
	requestInfo.UrlTemplate = u.UrlTemplate
	requestInfo.Method = abstractions.PUT
	requestInfo.SetStreamContent(content)
	return requestInfo
}
