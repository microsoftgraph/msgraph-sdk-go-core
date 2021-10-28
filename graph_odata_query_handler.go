package msgraphgocore

import (
	nethttp "net/http"
	regexp "regexp"

	abs "github.com/microsoft/kiota/abstractions/go"
	khttp "github.com/microsoft/kiota/http/go/nethttp"
)

// GraphODataQueryHandler is a handler that adds the dollar sign in front of OData query parameters
type GraphODataQueryHandler struct {
	regex          *regexp.Regexp
	handlerOptions GraphODataQueryHandlerOptions
}

type GraphODataQueryHandlerOptions struct {
	shouldReplace func(*nethttp.Request) bool
}

type graphODataQueryHandlerOptionsInt interface {
	abs.RequestOption
	ShouldReplace(*nethttp.Request) bool
}

func (o *GraphODataQueryHandlerOptions) ShouldReplace(req *nethttp.Request) bool {
	return o.shouldReplace(req)
}

var keyValue = abs.RequestOptionKey{
	Key: "GraphODataQueryHandler",
}

func (o *GraphODataQueryHandlerOptions) GetKey() abs.RequestOptionKey {
	return keyValue
}

// NewGraphODataQueryHandler creates a new instance of GraphODataQueryHandler
func NewGraphODataQueryHandler() *GraphODataQueryHandler {
	return NewGraphODataQueryHandlerWithOptions(
		GraphODataQueryHandlerOptions{
			shouldReplace: func(*nethttp.Request) bool {
				return true
			},
		})
}

// NewGraphODataQueryHandlerWithOptions creates a new instance of GraphODataQueryHandler
// Parameters:
// 		options: GraphODataQueryHandlerOptions options to use for the handler
func NewGraphODataQueryHandlerWithOptions(options GraphODataQueryHandlerOptions) *GraphODataQueryHandler {
	replacementRegexp := "(?i)([^$])(count|expand|filter|format|orderby|search|select|skip|skiptoken|top)="
	return &GraphODataQueryHandler{
		regex:          regexp.MustCompile(replacementRegexp),
		handlerOptions: options,
	}
}

func (middleware GraphODataQueryHandler) Intercept(pipeline khttp.Pipeline, req *nethttp.Request) (*nethttp.Response, error) {
	reqOption, ok := req.Context().Value(keyValue).(graphODataQueryHandlerOptionsInt)
	if ok && reqOption.ShouldReplace(req) || !ok && middleware.handlerOptions.shouldReplace(req) {
		req.URL.RawQuery = middleware.regex.ReplaceAllString("?"+req.URL.RawQuery, "$1$$$2=")[1:]
		// inserting and removing the ? sign so we can make no dollar mandatory and avoid adding a second dollar when already here
	}
	return pipeline.Next(req)
}
