package msgraphgocore

import (
	nethttp "net/http"
	regexp "regexp"

	khttp "github.com/microsoft/kiota/http/go/nethttp"
)

// GraphODataQueryHandler is a handler that adds the dollar sign in front of OData query parameters
type GraphODataQueryHandler struct {
	regex *regexp.Regexp
}

// NewGraphODataQueryHandler creates a new instance of GraphODataQueryHandler
func NewGraphODataQueryHandler() *GraphODataQueryHandler {
	//TODO constructor with options
	replacementRegexp := "(?i)([^$])(count|expand|filter|format|orderby|search|select|skip|skiptoken|top)="
	return &GraphODataQueryHandler{
		regex: regexp.MustCompile(replacementRegexp),
	}
}

func (middleware GraphODataQueryHandler) Intercept(pipeline khttp.Pipeline, req *nethttp.Request) (*nethttp.Response, error) {
	//TODO get options from context, default back to properties options
	req.URL.RawQuery = middleware.regex.ReplaceAllString("?"+req.URL.RawQuery, "$1$$$2=")[1:]
	// inserting and removing the ? sign so we can make no dollar mandatory and avoid adding a second dollar when already here
	return pipeline.Next(req)
}
