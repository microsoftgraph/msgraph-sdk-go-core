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
	ShouldReplace func(*nethttp.Request) bool
}

type graphODataQueryHandlerOptionsInt interface {
	abs.RequestOption
	GetShouldReplace() func(*nethttp.Request) bool
}

func (o *GraphODataQueryHandlerOptions) GetShouldReplace() func(req *nethttp.Request) bool {
	return o.ShouldReplace
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
			ShouldReplace: func(*nethttp.Request) bool {
				return true
			},
		})
}

// NewGraphODataQueryHandlerWithOptions creates a new instance of GraphODataQueryHandler
func NewGraphODataQueryHandlerWithOptions(options GraphODataQueryHandlerOptions) *GraphODataQueryHandler {
	replacementRegexp := "(?i)([^$|4])(count|deltatoken|expand|filter|format|orderby|search|select|skip|skiptoken|top)="
	return &GraphODataQueryHandler{
		regex:          regexp.MustCompile(replacementRegexp),
		handlerOptions: options,
	}
}

func (middleware GraphODataQueryHandler) Intercept(pipeline khttp.Pipeline, middlewareIndex int, req *nethttp.Request) (*nethttp.Response, error) {
	reqOption, ok := req.Context().Value(keyValue).(graphODataQueryHandlerOptionsInt)
	if ok && reqOption.GetShouldReplace()(req) || !ok && middleware.handlerOptions.ShouldReplace(req) {
		req.URL.RawQuery = middleware.regex.ReplaceAllString("?"+req.URL.RawQuery, "$1$$$2=")[1:]
		// inserting and removing the ? sign so we can make no dollar mandatory and avoid adding a second dollar when already here
	}
	return pipeline.Next(req, middlewareIndex)
}
