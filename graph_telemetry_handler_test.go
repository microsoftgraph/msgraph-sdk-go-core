package msgraphgocore

import (
	nethttp "net/http"
	httptest "net/http/httptest"
	testing "testing"

	khttp "github.com/microsoft/kiota/http/go/nethttp"
	assert "github.com/stretchr/testify/assert"
)

type NoopMiddleware struct {
	response *nethttp.Response
}

func (m *NoopMiddleware) Do(req *nethttp.Request) (*nethttp.Response, error) {
	return m.response, nil
}
func (m *NoopMiddleware) GetNext() khttp.Middleware {
	return nil
}
func (m *NoopMiddleware) SetNext(value khttp.Middleware) {

}

func TestItCreatesANewHandler(t *testing.T) {
	handler := NewGraphTelemetryHandler(&GraphClientOptions{})
	if handler == nil {
		t.Error("handler is nil")
	}
}

func TestItAddsHeaders(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(res nethttp.ResponseWriter, req *nethttp.Request) {
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))
	defer func() { testServer.Close() }()
	handler := NewGraphTelemetryHandler(&GraphClientOptions{})
	handler.SetNext(&NoopMiddleware{})
	req, err := nethttp.NewRequest(nethttp.MethodGet, testServer.URL, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := handler.Do(req)
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, resp)
	sdkVersionHeaderValue := req.Header[nethttp.CanonicalHeaderKey("SdkVersion")]
	assert.NotEmpty(t, sdkVersionHeaderValue)
	assert.Contains(t, sdkVersionHeaderValue[0], "graph-go-core")
	assert.Contains(t, sdkVersionHeaderValue[0], "hostOS")
	assert.Contains(t, sdkVersionHeaderValue[0], "hostArch")
	assert.Contains(t, sdkVersionHeaderValue[0], "runtimeEnvironment")
	assert.NotEmpty(t, req.Header[nethttp.CanonicalHeaderKey("client-request-id")])
}

func TestItAddsServiceLibInfo(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(res nethttp.ResponseWriter, req *nethttp.Request) {
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))
	defer func() { testServer.Close() }()
	handler := NewGraphTelemetryHandler(&GraphClientOptions{
		GraphServiceLibraryVersion: "1.0.0",
	})
	handler.SetNext(&NoopMiddleware{})
	req, err := nethttp.NewRequest(nethttp.MethodGet, testServer.URL, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := handler.Do(req)
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, resp)
	sdkVersionHeaderValue := req.Header[nethttp.CanonicalHeaderKey("SdkVersion")]
	assert.NotEmpty(t, sdkVersionHeaderValue)
	assert.Contains(t, sdkVersionHeaderValue[0], "graph-go/")
}

func TestItAddsServiceInfo(t *testing.T) {
	testServer := httptest.NewServer(nethttp.HandlerFunc(func(res nethttp.ResponseWriter, req *nethttp.Request) {
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))
	defer func() { testServer.Close() }()
	handler := NewGraphTelemetryHandler(&GraphClientOptions{
		GraphServiceLibraryVersion: "1.0.0",
		GraphServiceVersion:        "v1",
	})
	handler.SetNext(&NoopMiddleware{})
	req, err := nethttp.NewRequest(nethttp.MethodGet, testServer.URL, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := handler.Do(req)
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, resp)
	sdkVersionHeaderValue := req.Header[nethttp.CanonicalHeaderKey("SdkVersion")]
	assert.NotEmpty(t, sdkVersionHeaderValue)
	assert.Contains(t, sdkVersionHeaderValue[0], "graph-go-v1/")
}
