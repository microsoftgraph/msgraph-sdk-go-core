package msgraphgocore

import (
	nethttp "net/http"
	httptest "net/http/httptest"
	testing "testing"

	assert "github.com/stretchr/testify/assert"
)

type NoopPipeline struct {
	client *nethttp.Client
}

func (pipeline *NoopPipeline) Next(req *nethttp.Request) (*nethttp.Response, error) {
	return pipeline.client.Do(req)
}
func newNoopPipeline() *NoopPipeline {
	return &NoopPipeline{
		client: nethttp.DefaultClient,
	}
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
	req, err := nethttp.NewRequest(nethttp.MethodGet, testServer.URL, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := handler.Intercept(newNoopPipeline(), req)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, resp)
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
	req, err := nethttp.NewRequest(nethttp.MethodGet, testServer.URL, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := handler.Intercept(newNoopPipeline(), req)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, resp)
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
	req, err := nethttp.NewRequest(nethttp.MethodGet, testServer.URL, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := handler.Intercept(newNoopPipeline(), req)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, resp)
	sdkVersionHeaderValue := req.Header[nethttp.CanonicalHeaderKey("SdkVersion")]
	assert.NotEmpty(t, sdkVersionHeaderValue)
	assert.Contains(t, sdkVersionHeaderValue[0], "graph-go-v1/")
}
