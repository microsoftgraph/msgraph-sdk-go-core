module github.com/microsoftgraph/msgraph-sdk-go-core

go 1.17

require (
	github.com/google/uuid v1.3.0
	github.com/microsoft/kiota-abstractions-go v0.3.0
	github.com/microsoft/kiota-http-go v0.1.0
	github.com/microsoft/kiota-serialization-json-go v0.2.0
	github.com/stretchr/testify v1.7.1
)

require (
	github.com/cjlapao/common-go v0.0.19 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

retract (
	v0.11.0
	// error in version bump, bumped minor instead of patch, causing issues with update commands as long as we don't have a higher version number
	v0.0.14
// contains retraction only
)
