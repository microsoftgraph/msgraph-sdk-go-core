package msgraphgocore

import (
	nethttp "net/http"

	runtime "runtime"

	uuid "github.com/google/uuid"
	kmiddleware "github.com/microsoft/kiota/http/go/nethttp/middleware"
)

// GraphTelemetryHandler is a middleware handler that adds telemetry headers to requests.
type GraphTelemetryHandler struct {
	kmiddleware.CallbackHandler
}

// NewGraphTelemetryHandler creates a new GraphTelemetryHandler.
// Parameters:
//   options - the options for the GraphClient.
// Returns:
//   the new GraphTelemetryHandler.
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
		featuresSuffix := ""
		if runtime.GOOS != "" {
			featuresSuffix += " hostOS=" + runtime.GOOS + ";"
		}
		if runtime.GOARCH != "" {
			featuresSuffix += " hostArch=" + runtime.GOARCH + ";"
		}
		goVersion := runtime.Version()
		if goVersion != "" {
			featuresSuffix += " runtimeEnvironment=" + goVersion + ";"
		}
		if featuresSuffix != "" {
			featuresSuffix = " (" + featuresSuffix[1:] + ")"
		}
		req.Header.Add("SdkVersion", "graph-go-core/"+CoreVersion+serviceVersionSuffix+featuresSuffix)
		req.Header.Add("client-request-id", uuid.NewString())
		return nil
	}
	return &GraphTelemetryHandler{
		CallbackHandler: *kmiddleware.NewTelemetryHandler(callback, nil),
	}
}
