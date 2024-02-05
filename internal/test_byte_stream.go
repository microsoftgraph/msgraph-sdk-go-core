package internal

import (
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"io"
	"io/fs"
	"os"
	"time"
)

type MockByteStream struct {
	Content []byte
	offset  int
}

func (m *MockByteStream) Stat() (os.FileInfo, error) {
	return &fakeFileInfo{
		dir:      false,
		basename: "mockByteStream",
		modtime:  time.Time{},
		ents:     nil,
		contents: string(m.Content),
		err:      nil,
	}, nil
}

func (m *MockByteStream) Read(p []byte) (n int, err error) {
	if m.offset >= len(m.Content) {
		return 0, io.EOF
	}
	v := copy(p, m.Content[m.offset:])
	m.offset += v
	return v, nil
}

func (m *MockByteStream) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

type fakeFileInfo struct {
	dir      bool
	basename string
	modtime  time.Time
	ents     []*fakeFileInfo
	contents string
	err      error
}

func (f *fakeFileInfo) Name() string       { return f.basename }
func (f *fakeFileInfo) Sys() any           { return nil }
func (f *fakeFileInfo) ModTime() time.Time { return f.modtime }
func (f *fakeFileInfo) IsDir() bool        { return f.dir }
func (f *fakeFileInfo) Size() int64        { return int64(len(f.contents)) }
func (f *fakeFileInfo) Mode() fs.FileMode {
	if f.dir {
		return 0755 | fs.ModeDir
	}
	return 0644
}

type UploadResponse struct {
	// Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
	additionalData map[string]interface{}
}

func (s *UploadResponse) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

func (s *UploadResponse) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return make(map[string]func(serialization.ParseNode) error)
}

func CreateUploadResponseFromDiscriminatorValue(parseNode serialization.ParseNode) (serialization.Parsable, error) {
	res := UploadResponse{}
	return &res, nil
}

type UploadResponseble interface {
	Serialize(writer serialization.SerializationWriter) error
	GetFieldDeserializers() map[string]func(serialization.ParseNode) error
}
