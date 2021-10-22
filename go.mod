module github.com/microsoftgraph/msgraph-sdk-go-core

replace github.com/microsoft/kiota/abstractions/go => ./kiota/abstractions/go

replace github.com/microsoft/kiota/serialization/go/json => ./kiota/serialization/go/json

replace github.com/microsoft/kiota/http/go/nethttp => ./kiota/http/go/nethttp

replace github.com/microsoft/kiota/authentication/go/azure => ./kiota/authentication/go/azure

//TODO update references and remove replaces + git submodule once the libraries are on their own repo

go 1.16

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/microsoft/kiota/abstractions/go v0.0.0-20211020104304-4deb4d4c4659 // indirect
	github.com/microsoft/kiota/authentication/go/azure v0.0.0-00010101000000-000000000000 // indirect
	github.com/microsoft/kiota/http/go/nethttp v0.0.0-00010101000000-000000000000 // indirect
	github.com/microsoft/kiota/serialization/go/json v0.0.0-20211020104304-4deb4d4c4659 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
)
