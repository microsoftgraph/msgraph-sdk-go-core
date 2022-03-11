module github.com/microsoftgraph/msgraph-sdk-go-core

go 1.17

require (
	github.com/google/uuid v1.3.0
	github.com/microsoft/kiota/abstractions/go v0.0.0-20220304192020-9b3bd245842e
	github.com/microsoft/kiota/http/go/nethttp v0.0.0-20220303111159-de55f78b58f1
	github.com/microsoft/kiota/serialization/go/json v0.0.0-20220304192020-9b3bd245842e
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.21.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.9.1 // indirect
	github.com/cjlapao/common-go v0.0.18 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/microsoft/kiota/authentication/go/azure v0.0.0-20220203083330-462603bf370f // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.1 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

retract (
	v0.11.0
	// error in version bump, bumped minor instead of patch, causing issues with update commands as long as we don't have a higher version number
	v0.0.14
	// contains retraction only
)
