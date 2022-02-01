package msgraphgocore

import (
	"net/http"
	nethttp "net/http"
	"net/http/httptrace"

	khttp "github.com/microsoft/kiota/http/go/nethttp"
)

// LoggingHandler represents a middleware used to print logs about http requests
type LoggingHandler struct {
	trace *httptrace.ClientTrace
}

// NewLoggingHandler creates an instance of the logging handler middleware
func NewLoggingHandler(clientTrace *httptrace.ClientTrace) *LoggingHandler {
	return &LoggingHandler{
		trace: clientTrace,
	}
}

func (logger LoggingHandler) Intercept(pipeline khttp.Pipeline, middlewareIndex int, req *nethttp.Request) (*http.Response, error) {
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), logger.trace))
	return pipeline.Next(req, middlewareIndex)
}
