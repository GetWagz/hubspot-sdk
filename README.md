# Hubspot SDK for Go

[![GoDoc](https://godoc.org/github.com/getwagz/hubspot-sdk?status.svg)](https://godoc.org/github.com/getwagz/hubspot-sdk)

This library serves as a simple SDK for the Hubspot API. The existing solutions did not meet our needs, so we decided to roll our own. We do not currently implement the entire API (as that is a lot) so PRs are always welcome.

*NOTE* This SDK exists as a temporary solution until an official SDK is released from Hubspot or a community-standard SDK is developed. 

## Installing

You can simply install the package:

`go get github.com/getwagz/hubspot-sdk`

Or if you are using `dep`:

`dep ensure -add github.com/getwagz/hubspot-sdk`

## Usage

First, there are some optional environment variables (with *hopefully* sane defaults):

`HUBSPOT_SDK_ENV` is the environment the SDK is running in; defaults to `development`

`HUBSPOT_SDK_ROOT_URL` is the root URL for the Hubspot service, defaults to `https://api.hubapi.com`

`HUBSPOT_SDK_APPLICATION_ID` is the oAuth application ID; this is not required for all functionality

`HUBSPOT_SDK_USER_ID` is the oAuth userID for the developer account; this is not required for all functionality

`HUBSPOT_SDK_API_KEY` is the api key for your account; defaults to `demo`

`HUBSPOT_SDK_LOGGING` will toggle logging on if not blank and not `off`, `false`, or `no`

## Testing

Please note that testing without changing environment variables will only be able to test some aspects of the API; the `demo` api key for the `HUBSPOT_SDK_API_KEY` will allow testing `Contacts` but not `Events`, which require an oAuth application. If testing those is important to you, you should use a dummy account and pass in the information as appropriate in the environment.

To run the tests, run

`go test`

For coverage in HTML format, run

`go test -coverprofile=coverage.out && go tool cover -html=coverage.out`

The coverage is notably lower than ideal, which may cause concerns. However, most of the uncovered calls would be calls directly to the Hubspot API, which we cannot easily mock in success conditions, that are malformed. Feel free to check the results of the coverage report to see what exactly isn't covered and make a determination if that is acceptable to you. This library is currently being used in production.

## Contributing

Pull Requests are welcome! See our `CONTRIBUTING.md` file for more information.

## Third-party Libraries

The following libraries are used in this project. We thank the creators and maintainers for making our lives easier!

[Resty](https://github.com/go-resty/resty)

[Logrus](https://github.com/sirupsen/logrus)

[Testify](https://github.com/stretchr/testify)

[Mapstructure](https://github.com/mitchellh/mapstructure)

## Endpoints Implemented

- Contacts
  - Create or Update [Docs](https://developers.hubspot.com/docs/methods/contacts/create_or_update)

## TODO

- Clean up any TODO: in the comments of the source
- Convert to Go Modules
- Improve Documentation
- Add more end points