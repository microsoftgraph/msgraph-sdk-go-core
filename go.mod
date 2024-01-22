module github.com/microsoftgraph/msgraph-sdk-go-core

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.9.1
	github.com/google/uuid v1.5.0
	github.com/microsoft/kiota-abstractions-go v1.5.6
	github.com/microsoft/kiota-authentication-azure-go v1.0.1
	github.com/microsoft/kiota-http-go v1.1.1
	github.com/microsoft/kiota-serialization-json-go v1.0.5
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.5.1 // indirect
	github.com/cjlapao/common-go v0.0.39 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/std-uritemplate/std-uritemplate/go v0.0.50 // indirect
	go.opentelemetry.io/otel v1.22.0 // indirect
	go.opentelemetry.io/otel/metric v1.22.0 // indirect
	go.opentelemetry.io/otel/trace v1.22.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v0.11.0
	// error in version bump, bumped minor instead of patch, causing issues with update commands as long as we don't have a higher version number
	v0.0.14
// contains retraction only
)
