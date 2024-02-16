package fileuploader

import (
	"strings"
	"time"
)

type rangePair struct {
	Start float64
	End   float64
}

func stringIsNullOrEmpty(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || len(strings.TrimSpace(s)) == 0 {
		return false
	}
	return true
}

type UploadSession interface {
	GetExpirationDateTime() *time.Time
	GetNextExpectedRanges() []string
	GetOdataType() *string
	GetUploadUrl() *string
}

type ProgressCallBack func(current float64, total float64)

type UploadResult[T interface{}] interface {
	SetItemResponse(response T)
	GetItemResponse() T
	SetUploadSession(uploadSession UploadSession)
	GetUploadSession() UploadSession
	SetURI(uri string)
	GetURI() string
	SetUploadSucceeded(isSuccessful bool)
	GetUploadSucceeded() bool
}

func NewUploadResult[T interface{}]() UploadResult[T] {
	return &uploadResult[T]{}
}

type uploadResult[T interface{}] struct {
	itemResponse    T
	uploadSession   UploadSession
	uri             string
	uploadSucceeded bool
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

func (u *uploadResult[T]) SetURI(uri string) {
	u.uri = uri
}

func (u *uploadResult[T]) GetURI() string {
	return u.uri
}

func (u *uploadResult[T]) SetUploadSucceeded(isSuccessful bool) {
	u.uploadSucceeded = isSuccessful
}

func (u *uploadResult[T]) GetUploadSucceeded() bool {
	return u.uploadSucceeded
}
