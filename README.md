# Hubspot SDK for Go

[![GoDoc](https://godoc.org/github.com/getwagz/hubspot-sdk?status.svg)](https://godoc.org/github.com/getwagz/hubspot-sdk)

This library serves as a simple SDK for the Hubspot API. The existing solutions did not meet our needs, so we decided to roll our own. We do not currently implement the entire API (as that is a lot) so PRs are always welcome.

## Warning

This SDK is not yet ready for production use and is under active development. If you would like to assist, feel free, but know that it is actively changing and should *not* currently be used for any real application.

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

`HUBSPOT_SDK_LOGGING` will toggle logging on if not blank and not `off`, `false`, or `no`

`HUBSPOT_SDK_API_KEY` is the api key for your account; defaults to `demo`

`HUBSPOT_SDK_APPLICATION_ID` is the oAuth application ID; this is not required for all functionality but is required for complete non-mocked testing and calls requiring oAuth

`HUBSPOT_SDK_USER_ID` is the oAuth userID for the developer account; this is not required for all functionality but is required for complete non-mocked testing and calls requiring oAuth

`HUBSPOT_SDK_CLIENT_ID` is the oAuth client ID for the developer account; this is not required for all functionality but is required for complete non-mocked testing and calls requiring oAuth

`HUBSPOT_SDK_CLIENT_SECRET` is the oAuth client secret key for the developer account; this is not required for all functionality but is required for complete non-mocked testing and calls requiring oAuth

`HUBSPOT_SDK_OAUTH_REFRESH_TOKEN` is the oAuth refresh token for the developer account; this is not required for all functionality but is required for complete non-mocked testing and calls requiring oAuth

### OAuth

Many calls (for example, `events` on the `timeline`) require an oAuth token to complete. This requires setting up an account on the Hubspot Developer portal and then providing the correct environment variables. On `init()`, a request is made to get a new access token based upon the refresh token. If this fails, a log with level `error` is raised. *HOWEVER* we will not panic just because we cannot talk to the Hubspot API Server. We will attempt to refresh the token upon expiry.

*Given* that the API calls requiring oAuth will be mocked if these tokens are not provided, it is **VERY** important that you consider testing this API with the proper oAuth tokens to ensure the communication is correct. We are not responsible for a failure to test in your specific environment!

## Testing

Please note that testing without changing environment variables will only be able to test some aspects of the API; the `demo` api key for the `HUBSPOT_SDK_API_KEY` will allow testing `Contacts` but not `Events`, which require an oAuth application. If testing those is important to you, you should use a dummy account and pass in the information as appropriate in the environment.

To run the tests, run

`go test`

For coverage in HTML format, run

`go test -coverprofile=coverage.out && go tool cover -html=coverage.out`

The coverage is notably lower than ideal, which may cause concerns. However, most of the uncovered calls would be calls directly to the Hubspot API, which we cannot easily mock in success conditions, that are malformed. Feel free to check the results of the coverage report to see what exactly isn't covered and make a determination if that is acceptable to you.

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
  - Delete [Doc](https://developers.hubspot.com/docs/methods/contacts/delete_contact)
  - Get by Email [Doc](https://developers.hubspot.com/docs/methods/contacts/get_contact_by_email)
- Events
  - Create Event Type [Doc](https://developers.hubspot.com/docs/methods/timeline/create-event-type)
  - Delete Event Type [Doc](https://developers.hubspot.com/docs/methods/timeline/delete-event-type)
  - Create Event on Timeline [Doc](https://developers.hubspot.com/docs/methods/timeline/create-or-update-event)

## TODO

- Clean up any TODO: in the comments of the source
- Convert to Go Modules
- Improve Documentation
- Add more end points
  - Engagements [Doc](https://developers.hubspot.com/docs/methods/engagements/create_engagement)
