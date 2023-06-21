package fileuploader

import (
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"math"
	"strconv"
	"strings"
)

type LargeFileUploadTask[T interface{}] interface {
	UploadAsync(progress ProgressCallBack) UploadResult[T]
}

type largeFileUploadTask[T interface{}] struct {
	uploadSession   UploadSession
	adapter         abstractions.RequestAdapter
	fileContent     []byte
	maxSlice        float64
	parsableFactory serialization.ParsableFactory
}

func NewLargeFileUploadTask[T interface{}](adapter abstractions.RequestAdapter, uploadSession UploadSession, fileContent []byte, maxSlice float64, parsableFactory serialization.ParsableFactory) LargeFileUploadTask[T] {
	return &largeFileUploadTask[T]{
		adapter:         adapter,
		uploadSession:   uploadSession,
		fileContent:     fileContent,
		maxSlice:        maxSlice,
		parsableFactory: parsableFactory,
	}
}

// UploadAsync TODO Update function to use go routines
// TODO allow re-uploading slices to a maximum of 10
func (l *largeFileUploadTask[T]) UploadAsync(progress ProgressCallBack) UploadResult[T] {
	/*maxTries := 3
	uploadTries := 0

	for uploadTries <= maxTries {
		fmt.Println(uploadTries)
		uploadTries++
	}*/

	slices := l.createUploadSlices()
	for _, slice := range slices {
		_, _ = slice.UploadAsync() // check if successful
		progress(slice.RangeEnd, slice.TotalSessionLength)
	}
	panic("implement me")
}

func (l *largeFileUploadTask[T]) getRangesRemaining() []rangePair {
	rangePairs := make([]rangePair, len(l.uploadSession.GetNextExpectedRanges()))

	for i, ranges := range l.uploadSession.GetNextExpectedRanges() {
		rangeValues := strings.Split(ranges, "-")

		var startRange float64
		if s, err := strconv.ParseFloat(rangeValues[0], 64); err == nil {
			startRange = s
		}

		var endRange float64
		if !stringIsNullOrEmpty(rangeValues[1]) {
			if s, err := strconv.ParseFloat(rangeValues[1], 64); err == nil {
				endRange = s
			}
		} else {
			endRange = float64(len(l.fileContent))
		}

		rangePairs[i] = rangePair{
			Start: startRange,
			End:   endRange,
		}
	}

	return rangePairs
}

func (l largeFileUploadTask[T]) nextSliceLength(rangeBegin float64, rangeEnd float64) float64 {
	sizeBasedOnRange := rangeEnd - rangeBegin + 1
	return math.Min(sizeBasedOnRange, l.maxSlice)
}
