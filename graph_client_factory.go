package msgraphgocore

import (
	nethttp "net/http"

	khttp "github.com/microsoft/kiota/http/go/nethttp"
)

// Creates a new default set of middlewares for the Graph Client
// Parameters:
// 		options - the options to use for the middlewares
// Returns:
// 		the middlewares
func GetDefaultMiddlewaresWithOptions(options *GraphClientOptions) []khttp.Middleware {
	kiotaMiddlewares := khttp.GetDefaultMiddlewares()
	graphMiddlewares := []khttp.Middleware{
		NewGraphTelemetryHandler(options),
		NewGraphODataQueryHandler(),
	}
	graphMiddlewaresLen := len(graphMiddlewares)
	resultMiddlewares := make([]khttp.Middleware, len(kiotaMiddlewares)+graphMiddlewaresLen)
	copy(resultMiddlewares, graphMiddlewares)
	copy(resultMiddlewares[graphMiddlewaresLen:], kiotaMiddlewares)
	return resultMiddlewares
}

// Create a new default net/http client with the options configured for the Graph Client
// Parameters:
// 		middleware - the middlewares to use for the client
// Returns:
// 		the client
func GetDefaultClient(middleware ...khttp.Middleware) *nethttp.Client {
	return khttp.GetDefaultClient(middleware...)
}
