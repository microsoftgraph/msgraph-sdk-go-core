package fileuploader

import (
	"context"
	"fmt"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"io"
	"time"
)

const binaryContentType = "application/octet-steam"

type uploadSlice[T serialization.Parsable] struct {
	RequestAdapter     abstractions.RequestAdapter
	UrlTemplate        string
	RangeBegin         int64
	RangeEnd           int64
	TotalSessionLength int64
	RangeLength        int64
	byteStream         ByteStream
	errorMappings      abstractions.ErrorMappings
}

func (l *largeFileUploadTask[T]) createUploadSlices() []uploadSlice[T] {

	requestRanges := l.getRangesRemaining()
	maxSlice := l.maxSlice
	totalSessionLength := l.fileSize()

	// compute the correct upload ranges by splitting the values of ranges remaining from start to end
	var uploadSlices []uploadSlice[T]
	for _, v := range requestRanges {
		start := v.Start
		for start < totalSessionLength && start <= v.End {
			end := minOf(v.End, (start+maxSlice)-1, totalSessionLength-1)
			uploadSlices = append(uploadSlices, uploadSlice[T]{
				RequestAdapter:     l.adapter,
				UrlTemplate:        *l.uploadSession.GetUploadUrl(),
				RangeBegin:         start,
				RangeEnd:           end,
				RangeLength:        end - start + 1,
				TotalSessionLength: totalSessionLength,
				errorMappings:      l.errorMappings,
				byteStream:         l.byteStream,
			})
			start = end + 1
		}
	}

	return uploadSlices
}

func (l *largeFileUploadTask[T]) readSection(start, end int64) ([]byte, error) {
	// Calculate the length of the section to read
	length := (end - start) + 1

	// Seek to the start position
	_, err := l.byteStream.Seek(start, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Read the section into a buffer
	buffer := make([]byte, length)
	_, err = io.ReadFull(l.byteStream, buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func minOf(vars ...int64) int64 {
	minimum := vars[0]
	for _, i := range vars {
		if minimum > i {
			minimum = i
		}
	}
	return minimum
}

func (u *uploadSlice[T]) Upload(parsableFactory serialization.ParsableFactory) (interface{}, error) {
	data, err := u.readSection(u.RangeBegin, u.RangeEnd)
	if err != nil {
		return nil, err
	}
	requestInfo := u.createRequestInformation(data)

	// limit the upload time per slice to 5 minutes
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	return u.RequestAdapter.Send(ctx, requestInfo, parsableFactory, u.errorMappings)
}

func (u *uploadSlice[T]) readSection(start, end int64) ([]byte, error) {
	// Calculate the length of the section to read
	length := (end - start) + 1

	// Seek to the start position
	_, err := u.byteStream.Seek(start, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Read the section into a buffer
	buffer := make([]byte, length)
	_, err = io.ReadFull(u.byteStream, buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func (u *uploadSlice[T]) createRequestInformation(content []byte) *abstractions.RequestInformation {
	headers := abstractions.NewRequestHeaders()
	headers.Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", u.RangeBegin, u.RangeEnd, u.TotalSessionLength))
	headers.Add("Content-Length", fmt.Sprintf("%d", u.RangeLength))

	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Headers = headers
	requestInfo.UrlTemplate = u.UrlTemplate
	requestInfo.Method = abstractions.PUT
	requestInfo.SetStreamContentAndContentType(content, binaryContentType)
	return requestInfo
}
