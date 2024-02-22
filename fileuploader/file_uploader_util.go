package fileuploader

import (
	"strings"
	"time"
)

type rangePair struct {
	Start int64
	End   int64
}

func stringIsNullOrEmpty(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || len(s) == 0 {
		return true
	}
	return false
}

type UploadSession interface {
	GetExpirationDateTime() *time.Time
	SetExpirationDateTime(expirationDateTime *time.Time)
	GetNextExpectedRanges() []string
	SetNextExpectedRanges(nextExpectedRanges []string)
	GetOdataType() *string
	GetUploadUrl() *string
}

type ProgressCallBack func(current int64, total int64)

type UploadResult[T interface{}] interface {
	SetItemResponse(response T)
	GetItemResponse() T
	SetUploadSession(uploadSession UploadSession)
	GetUploadSession() UploadSession
	SetURI(uri *string)
	GetURI() *string
	SetUploadSucceeded(isSuccessful bool)
	GetUploadSucceeded() bool
	SetResponseErrors(errors []error)
	GetResponseErrors() []error
}

func NewUploadResult[T interface{}]() UploadResult[T] {
	return &uploadResult[T]{}
}

type uploadResult[T interface{}] struct {
	itemResponse    T
	uploadSession   UploadSession
	uri             *string
	uploadSucceeded bool
	responseErrors  []error
}

func (u *uploadResult[T]) SetItemResponse(response T) {
	u.itemResponse = response
}

func (u *uploadResult[T]) GetItemResponse() T {
	return u.itemResponse
}

func (u *uploadResult[T]) SetUploadSession(uploadSession UploadSession) {
	u.uploadSession = uploadSession
}

func (u *uploadResult[T]) GetUploadSession() UploadSession {
	return u.uploadSession
}

func (u *uploadResult[T]) SetURI(uri *string) {
	u.uri = uri
}

func (u *uploadResult[T]) GetURI() *string {
	return u.uri
}

func (u *uploadResult[T]) SetUploadSucceeded(isSuccessful bool) {
	u.uploadSucceeded = isSuccessful
}

func (u *uploadResult[T]) GetUploadSucceeded() bool {
	return u.uploadSucceeded
}

func (u *uploadResult[T]) SetResponseErrors(errors []error) {
	u.responseErrors = errors
}

func (u *uploadResult[T]) GetResponseErrors() []error {
	return u.responseErrors
}
