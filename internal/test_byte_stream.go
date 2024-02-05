package internal

import "os"

type MockByteStream struct {
	Content []byte
	offset  int
}

func (m *MockByteStream) Stat() (os.FileInfo, error) {
	var info os.FileInfo
	return info, nil
}

func (m *MockByteStream) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (m *MockByteStream) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (m *MockByteStream) Write(p []byte) (n int, err error) {
	return 0, nil
}
