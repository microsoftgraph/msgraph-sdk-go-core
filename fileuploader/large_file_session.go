package fileuploader

import (
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"time"
)

type UploadSessionResponse interface {
	serialization.Parsable
	GetExpirationDateTime() *time.Time
	SetExpirationDateTime(expirationDateTime *time.Time)
	GetNextExpectedRanges() []string
	SetNextExpectedRanges(nextExpectedRanges []string)
}

type largeFileUploadSession struct {
	expirationDateTime *time.Time
	nextExpectedRanges []string
}

func (l *largeFileUploadSession) Serialize(writer serialization.SerializationWriter) error {
	if l.expirationDateTime != nil {
		if err := writer.WriteTimeValue("expirationDateTime", l.expirationDateTime); err != nil {
			return err
		}
	}
	if l.nextExpectedRanges != nil {
		if err := writer.WriteCollectionOfStringValues("nextExpectedRanges", l.nextExpectedRanges); err != nil {
			return err
		}
	}
	return nil
}

func (l *largeFileUploadSession) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return map[string]func(serialization.ParseNode) error{
		"expirationDateTime": func(n serialization.ParseNode) error {
			val, err := n.GetTimeValue()
			if err != nil {
				return err
			}
			if val != nil {
				l.SetExpirationDateTime(val)
			}
			return nil
		},
		"nextExpectedRanges": func(n serialization.ParseNode) error {
			val, err := n.GetCollectionOfPrimitiveValues("string")
			if err != nil {
				return err
			}
			if val != nil {
				res := make([]string, len(val))
				for i, v := range val {
					if v != nil {
						res[i] = *(v.(*string))
					}
				}
				l.SetNextExpectedRanges(res)
			}
			return nil
		},
	}
}

func (l *largeFileUploadSession) GetExpirationDateTime() *time.Time {
	return l.expirationDateTime
}

func (l *largeFileUploadSession) SetExpirationDateTime(expirationDateTime *time.Time) {
	l.expirationDateTime = expirationDateTime
}

func (l *largeFileUploadSession) GetNextExpectedRanges() []string {
	return l.nextExpectedRanges
}

func (l *largeFileUploadSession) SetNextExpectedRanges(nextExpectedRanges []string) {
	l.nextExpectedRanges = nextExpectedRanges
}

func newLargeFileUploadSession() UploadSessionResponse {
	return &largeFileUploadSession{}
}

// CreateUploadSessionDiscriminator creates a new instance of the appropriate class based on discriminator value
func CreateUploadSessionDiscriminator(serialization.ParseNode) (serialization.Parsable, error) {
	return newLargeFileUploadSession(), nil
}
