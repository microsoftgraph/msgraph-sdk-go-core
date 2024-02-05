package fileuploader

import (
	"context"
	"errors"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type LargeFileUploadTask[T interface{}] interface {
	UploadAsync(progress ProgressCallBack) UploadResult[T]
}

// ByteStream is an interface that represents a stream of bytes
type ByteStream interface {
	io.Seeker
	io.Reader
	Stat() (os.FileInfo, error)
}

type largeFileUploadTask[T interface{}] struct {
	uploadSession   UploadSession
	adapter         abstractions.RequestAdapter
	byteStream      ByteStream // *os.File by default implements ByteStream
	maxSlice        int64
	parsableFactory serialization.ParsableFactory
	errorMappings   abstractions.ErrorMappings
}

func NewLargeFileUploadTask[T interface{}](adapter abstractions.RequestAdapter, uploadSession UploadSession, byteStream ByteStream, maxSlice int64, parsableFactory serialization.ParsableFactory, errorMappings abstractions.ErrorMappings) LargeFileUploadTask[T] {
	return &largeFileUploadTask[T]{
		adapter:         adapter,
		uploadSession:   uploadSession,
		byteStream:      byteStream,
		maxSlice:        maxSlice,
		parsableFactory: parsableFactory,
		errorMappings:   errorMappings,
	}
}

// UploadAsync uploads the byteStream in slices and returns the result of the upload
func (l *largeFileUploadTask[T]) UploadAsync(progress ProgressCallBack) UploadResult[T] {
	result := NewUploadResult[T]()
	var wg sync.WaitGroup
	slices := l.createUploadSlices()
	for _, slice := range slices {
		wg.Add(1)
		uploadSlice := slice
		go func() {
			defer wg.Done()
			l.uploadAsync(progress, uploadSlice, result)
		}()
	}

	wg.Wait()
	return result
}

// Resume uploads the byteStream in slices and returns the result of the upload
func (l *largeFileUploadTask[T]) Resume(progress ProgressCallBack) (UploadResult[T], error) {
	// check if next expected ranges is empty

	if len(l.uploadSession.GetNextExpectedRanges()) == 0 {
		return nil, errors.New("UploadSession does not have next expected ranges")
	}
	return l.UploadAsync(progress), nil
}

// Cancel cancels the upload
func (l *largeFileUploadTask[T]) Cancel() error {
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(abstractions.DELETE, *l.uploadSession.GetUploadUrl(), make(map[string]string))
	err := l.adapter.SendNoContent(context.Background(), requestInfo, l.errorMappings)
	return err
}

func (l *largeFileUploadTask[T]) uploadAsync(progress ProgressCallBack, slice uploadSlice[T], result UploadResult[T]) {
	maxRetry := 3
	retry := 1
	for retry < maxRetry {
		// store the result of the upload
		response, err := slice.UploadAsync(l.uploadSession, l.parsableFactory) // check if successful
		if err != nil {
			// if not successful, try again
			if retry >= maxRetry {
				result.SetUploadSucceeded(false)
			}

		} else {
			result.SetUploadSucceeded(true)
			result.SetItemResponse(response.GetItemResponse())
			progress(slice.RangeEnd, slice.TotalSessionLength)
			break
		}

		retry++
	}
}

func (l *largeFileUploadTask[T]) getRangesRemaining() []rangePair {
	rangePairs := make([]rangePair, len(l.uploadSession.GetNextExpectedRanges()))

	for i, ranges := range l.uploadSession.GetNextExpectedRanges() {
		rangeValues := strings.Split(ranges, "-")

		var startRange int64
		if s, err := strconv.ParseInt(rangeValues[0], 10, 64); err == nil {
			startRange = s
		}

		var endRange int64
		if !stringIsNullOrEmpty(rangeValues[1]) {
			if s, err := strconv.ParseInt(rangeValues[1], 10, 64); err == nil {
				if endRange > l.fileSize() {
					endRange = l.fileSize() - 1
				} else {
					endRange = s
				}
			}
		} else {
			endRange = l.fileSize() - 1
		}

		rangePairs[i] = rangePair{
			Start: startRange,
			End:   endRange,
		}
	}

	return rangePairs
}

// returns the size of a byteStream
func (l *largeFileUploadTask[T]) fileSize() int64 {
	fileInfo, _ := l.byteStream.Stat()
	return fileInfo.Size()
}

func (l *largeFileUploadTask[T]) nextSliceLength(rangeBegin int64, rangeEnd int64) int64 {
	sizeBasedOnRange := rangeEnd - rangeBegin + 1
	return minOf(sizeBasedOnRange, l.maxSlice)
}
