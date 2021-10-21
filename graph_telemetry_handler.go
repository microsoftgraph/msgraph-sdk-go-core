package msgraphgocore

import (
	nethttp "net/http"

	uuid "github.com/google/uuid"
	kmiddleware "github.com/microsoft/kiota/http/go/nethttp/middleware"
)

type GraphTelemetryHandler struct {
	kmiddleware.TelemetryHandler
}

func NewGraphTelemetryHandler(options *GraphClientOptions) *GraphTelemetryHandler {
	callback := func(req *nethttp.Request) error {
		serviceVersionSuffix := ""
		if options != nil && options.GraphServiceLibraryVersion != "" {
			serviceVersionSuffix += ", graph-go"
			if options.GraphServiceVersion != "" {
				serviceVersionSuffix += "-" + options.GraphServiceVersion
			}
			serviceVersionSuffix += "/" + options.GraphServiceLibraryVersion
		}
		req.Header.Add("SdkVersion", "graph-go-core/"+CoreVersion+serviceVersionSuffix)
		req.Header.Add("client-request-id", uuid.NewString())
		return nil
	}
	return &GraphTelemetryHandler{
		TelemetryHandler: *kmiddleware.NewTelemetryHandler(callback),
	}
}
