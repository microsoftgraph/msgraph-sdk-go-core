package msgraphgocore

// Options for the GraphClient
type GraphClientOptions struct {
	// The version of the targeted service for telemetry (v1.0, beta)
	GraphServiceVersion string
	// The version of the service library for telemetry (1.2.3)
	GraphServiceLibraryVersion string
}
