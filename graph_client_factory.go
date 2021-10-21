package msgraphgocore

import (
	khttp "github.com/microsoft/kiota/http/go/nethttp"
)

func GetDefaultMiddlewaresWithOptions(options *GraphClientOptions) []khttp.Middleware {
	kiotaMiddlewares := khttp.GetDefaultMiddlewares()
	graphMiddlewares := []khttp.Middleware{
		NewGraphTelemetryHandler(options),
	}
	graphMiddlewaresLen := len(graphMiddlewares)
	resultMiddlewares := make([]khttp.Middleware, len(kiotaMiddlewares)+graphMiddlewaresLen)
	copy(resultMiddlewares, graphMiddlewares)
	copy(resultMiddlewares[graphMiddlewaresLen:], kiotaMiddlewares)
	return resultMiddlewares
}
