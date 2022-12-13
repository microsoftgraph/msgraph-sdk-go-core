module github.com/microsoftgraph/msgraph-sdk-go-core

go 1.18

require (
	github.com/google/uuid v1.3.0
	github.com/microsoft/kiota-abstractions-go v0.15.0
	github.com/microsoft/kiota-http-go v0.10.0
	github.com/microsoft/kiota-serialization-json-go v0.7.2
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/cjlapao/common-go v0.0.36 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.opentelemetry.io/otel v1.11.2 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v0.11.0
	// error in version bump, bumped minor instead of patch, causing issues with update commands as long as we don't have a higher version number
	v0.0.14
// contains retraction only
)
