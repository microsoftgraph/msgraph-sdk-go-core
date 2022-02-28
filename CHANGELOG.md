# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

## [0.0.10] - 2022-02-28

### Changed

- Fixed a bug where http client configuration would impact the default client configuration for other usages.

## [0.0.9] - 2022-02-16

### Added

- Added support for deserializing error responses (will return error)

### Changed

- Fixed a bug where response body compression would send empty bodies

## [0.0.8] - 2022-02-08

### Added

- Added support for request body compression (gzip)
- Added support for response body decompression (gzip)

### Changed

- Fixes a bug where resuming the page iterator wouldn't work
- Fixes a bug where OData query parameters would be added twice in some cases

## [0.0.7] - 2022-02-03

### Changed

- Updated references to Kiota packages to fix a [bug where the access token would never be attached to the request](https://github.com/microsoft/kiota/pull/1116). 

## [0.0.6] - 2022-02-02

### Added

- Adds missing delta token for OData query parameters dollar sign injection.
- Adds PageIterator task

## [0.0.5] - 2021-12-02

### Changed

- Fixes a bug where the middleware pipeline would run only on the first request of the client/adapter/http client.

## [0.0.4] - 2021-12-01

### Changed

- Adds the missing github.com/microsoft/kiota/authentication/go/azure dependency

## [0.0.3] - 2021-11-30

### Changed

- Updated dependencies and switched to Go 17.

## [0.0.2] - 2021-11-08

### Changed

- Updated kiota abstractions and http to provide support for setting the base URL

## [0.0.1] - 2021-10-22

### Added

- Initial release
