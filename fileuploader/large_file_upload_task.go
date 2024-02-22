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
	"time"
)

type LargeFileUploadTask[T serialization.Parsable] interface {
	Upload(progress ProgressCallBack) UploadResult[T]
	Resume(progress ProgressCallBack) (UploadResult[T], error)
	RefreshUploadStatus() error
	Cancel() error
}

// ByteStream is an interface that represents a stream of bytes
type ByteStream interface {
	io.ReaderAt
	Stat() (os.FileInfo, error)
}

type largeFileUploadTask[T serialization.Parsable] struct {
	uploadSession   UploadSession
	adapter         abstractions.RequestAdapter
	byteStream      ByteStream // *os.File by default implements ByteStream
	maxSlice        int64
	parsableFactory serialization.ParsableFactory
	errorMappings   abstractions.ErrorMappings
}

func NewLargeFileUploadTask[T serialization.Parsable](adapter abstractions.RequestAdapter, uploadSession UploadSession, byteStream ByteStream, maxSlice int64, parsableFactory serialization.ParsableFactory, errorMappings abstractions.ErrorMappings) LargeFileUploadTask[T] {
	return &largeFileUploadTask[T]{
		adapter:         adapter,
		uploadSession:   uploadSession,
		byteStream:      byteStream,
		maxSlice:        maxSlice,
		parsableFactory: parsableFactory,
		errorMappings:   errorMappings,
	}
}

// Upload uploads the byteStream in slices and returns the result of the upload
func (l *largeFileUploadTask[T]) Upload(progress ProgressCallBack) UploadResult[T] {
	result := NewUploadResult[T]()
	slices := l.createUploadSlices()
	maxRetriesPerRequest := 3

	// slices of errors
	var responseErrors []error
	var itemResponse T
	var location *string

	var wg sync.WaitGroup
	wg.Add(len(slices))

	for _, slice := range slices {
		uploadSlice := slice
		go func() {
			defer wg.Done()
			response, uploadLocation, err := l.uploadWithRetry(uploadSlice, maxRetriesPerRequest)
			if err != nil {
				responseErrors = append(responseErrors, err)
			} else {
				progress(uploadSlice.RangeEnd, uploadSlice.TotalSessionLength)
			}
			if response != nil {
				itemResponse = response.(T)
			}
			location = uploadLocation
		}()
	}

	wg.Wait()

	if len(responseErrors) > 0 {
		result.SetUploadSucceeded(false)
		result.SetResponseErrors(responseErrors)
	} else {
		result.SetUploadSucceeded(true)
		result.SetUploadSession(l.uploadSession)
		result.SetItemResponse(itemResponse)
		result.SetURI(location)
	}

	return result
}

// Resume uploads the byteStream in slices and returns the result of the upload
func (l *largeFileUploadTask[T]) Resume(progress ProgressCallBack) (UploadResult[T], error) {
	err := l.RefreshUploadStatus()
	if err != nil {
		return nil, err
	}

	if len(l.uploadSession.GetNextExpectedRanges()) == 0 {
		return nil, errors.New("UploadSession does not have next expected ranges")
	}

	if l.uploadSession.GetExpirationDateTime().Before(time.Now()) {
		return nil, errors.New("UploadSession has expired")
	}

	return l.Upload(progress), nil
}

func (l *largeFileUploadTask[T]) RefreshUploadStatus() error {
	requestInfo := abstractions.NewRequestInformation()
	requestInfo.UrlTemplate = *l.uploadSession.GetUploadUrl()
	requestInfo.Method = abstractions.GET
	requestInfo.Headers.TryAdd("Accept", "application/json")

	result, err := l.adapter.Send(context.Background(), requestInfo, CreateUploadSessionDiscriminator, l.errorMappings)
	if err != nil {
		return err
	}

	sessionResponse := result.(UploadSessionResponse)

	l.uploadSession.SetExpirationDateTime(sessionResponse.GetExpirationDateTime())
	l.uploadSession.SetNextExpectedRanges(sessionResponse.GetNextExpectedRanges())

	return nil
}

// Cancel cancels the upload
func (l *largeFileUploadTask[T]) Cancel() error {
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(abstractions.DELETE, *l.uploadSession.GetUploadUrl(), make(map[string]string))
	err := l.adapter.SendNoContent(context.Background(), requestInfo, l.errorMappings)
	return err
}

func (l *largeFileUploadTask[T]) uploadWithRetry(slice uploadSlice[T], maxRetry int) (interface{}, *string, error) {
	retry := 1
	var parseable interface{}
	var location *string
	var err error
	for retry < maxRetry {
		// store the result of the upload
		parseable, location, err = slice.Upload(l.parsableFactory) // check if successful
		if err != nil {
			if retry >= maxRetry {
				return nil, nil, err
			}
			// backoff before retrying
			time.Sleep(time.Duration(retry) * time.Second)
		}
		retry++
	}
	return parseable, location, err
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
