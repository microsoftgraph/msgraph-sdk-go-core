package msgraphgocore

import (
	"fmt"
	"net/http"
	httptest "net/http/httptest"
	"net/http/httptrace"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceCallbacksGetCalled(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "hello world")
	}))
	defer testServer.Close()

	gotConnWasCalled := false

	trace := httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			gotConnWasCalled = true
		},
	}

	logHandler := NewLoggingHandler(&trace)
	client := GetDefaultClient(nil, logHandler)
	client.Get("https://example.com")

	assert.Equal(t, gotConnWasCalled, true)
}
