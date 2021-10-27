package msgraphgocore

import (
	nethttp "net/http"
	regexp "regexp"

	kmiddleware "github.com/microsoft/kiota/http/go/nethttp/middleware"
)

// GraphODataQueryHandler is a handler that adds the dollar sign in front of OData query parameters
type GraphODataQueryHandler struct {
	kmiddleware.CallbackHandler
}

// NewGraphODataQueryHandler creates a new instance of GraphODataQueryHandler
func NewGraphODataQueryHandler() *GraphODataQueryHandler {
	replacementRegexp := "(?i)([^$])(count|expand|filter|format|orderby|search|select|skip|skiptoken|top)="
	re := regexp.MustCompile(replacementRegexp)
	callback := func(req *nethttp.Request) error {
		req.URL.RawQuery = re.ReplaceAllString("?"+req.URL.RawQuery, "$1$$$2=")[1:]
		// inserting and removing the ? sign so we can make no dollar mandatory and avoid adding a second dollar when already here
		return nil
	}

	return &GraphODataQueryHandler{
		CallbackHandler: *kmiddleware.NewCallbackHandler(callback, nil),
	}
}
