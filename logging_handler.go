package msgraphgocore

import (
	"net/http"
	nethttp "net/http"
	"net/http/httptrace"

	khttp "github.com/microsoft/kiota/http/go/nethttp"
)

type LoggingHandler struct {
	trace *httptrace.ClientTrace
}

func NewLoggingHandler(clientTrace *httptrace.ClientTrace) *LoggingHandler {
	return &LoggingHandler{
		trace: clientTrace,
	}
}

func (logger LoggingHandler) Intercept(pipeline khttp.Pipeline, middlewareIndex int, req *nethttp.Request) (*http.Response, error) {
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), logger.trace))
	return pipeline.Next(req, middlewareIndex)
}
