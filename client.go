package hubspot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty"
)

// APIReturn represents a successful API return value, with the HTTPCode representing the exact returned HTTP Code and the Body
// representing the returned body
type APIReturn struct {
	HTTPCode int
	Body     interface{}
}

// APIError is the error struct containing additional information regarding the error that occurred. The HTTPCode is the returned value from
// Hubspot OR a 400 if there was an error prior to calling Hubspot. The SystemCode is a constant and you can check systemCodes.go for more information.
//
type APIError struct {
	HTTPCode   int
	SystemCode string
	Message    string
	Body       *map[string]interface{}
}

func (e APIError) Error() string {
	return e.Message
}

// prepareCall is the main entry point into the underlying HTTP req	uest generation and is used to actually prepare
// calls to the API. Mocking is handled here as well.
func prepareCall(endpoint string, pathParams map[string]string, data interface{}) (ret *APIReturn, err error) {
	// first, find the endpoint
	info, infoFound := endpoints[endpoint]
	if !infoFound {
		return nil, APIError{
			HTTPCode:   http.StatusBadRequest,
			SystemCode: "request_error_endpoint_not_found",
			Message:    "Could not find that endpoint",
		}
	}

	// if the oauth is required but not provided, then we just return the mocked data
	if (info.RequireOAuth && Config.HubSpotOAuthRefreshToken == "") || (Config.HubspotApplicationID == "test" && info.MockGood != nil) {
		return &APIReturn{
			HTTPCode: http.StatusOK,
			Body:     info.MockGood,
		}, nil
	}

	parsedPath := info.Path
	// we need to loop over the pathParams and fo a string replace
	for k, v := range pathParams {
		parsedPath = strings.Replace(parsedPath, k, v, -1)
	}

	// we always replace application id if it is there
	parsedPath = strings.Replace(parsedPath, ":applicationID", Config.HubspotApplicationID, -1)

	return makeCall(info.Method, parsedPath, data, info.RequireOAuth)
}

// makeCall makes the call to the Hubspot API
func makeCall(httpMethod, endpoint string, data interface{}, requireOAuth bool) (ret *APIReturn, err error) {
	if strings.HasPrefix(endpoint, "/") {
		endpoint = endpoint[1:]
	}

	url := fmt.Sprintf("%s%s", Config.RootURL, endpoint)

	var response *resty.Response

	request := resty.R().
		SetHeader("Accept", "application/json")

	queryParams := map[string]string{}

	if requireOAuth && oAuthToken.AccessToken != "" {
		request.SetAuthToken(oAuthToken.AccessToken)
	} else {
		// if oauth is required, we do not send up the api key
		queryParams["hapikey"] = Config.HubspotAPIKey
	}
	if Config.HubspotUserID != "" {
		queryParams["userId"] = Config.HubspotUserID
	}

	log("info", "api_url", fmt.Sprintf("Calling URL: %s: %s", httpMethod, url), map[string]interface{}{
		"query": queryParams,
	})

	// Now, do what we need to do depending on the method
	var reqErr error

	switch httpMethod {
	case http.MethodGet:
		if data != nil {
			// merge the two data sets
			dataParsed, dataParsedOK := data.(map[string]string)
			if !dataParsedOK {
				return nil, APIError{
					HTTPCode:   http.StatusBadRequest,
					SystemCode: "request_error_bad_query_string",
					Message:    "GET requests must use a map[string]string{} for data",
				}
			}
			for k, v := range dataParsed {
				queryParams[k] = v
			}
		}
		response, reqErr = request.SetQueryParams(queryParams).Get(url)
	case http.MethodDelete:
		response, reqErr = request.SetQueryParams(queryParams).Delete(url)
	case http.MethodPost:
		response, reqErr = request.SetQueryParams(queryParams).SetBody(data).Post(url)
		fmt.Printf("\n%+v\n", request)
	case http.MethodPut:
		response, reqErr = request.SetQueryParams(queryParams).SetBody(data).Put(url)
	}

	if reqErr != nil {
		log("warning", "unknown_api_error", "we encountered an error calling the API", map[string]interface{}{
			"err": reqErr,
		})
		return nil, &APIError{
			HTTPCode:   http.StatusInternalServerError,
			SystemCode: "request_error",
			Message:    reqErr.Error(),
		}
	}

	if response == nil {
		// there is an unknown error from the server
		log("warning", "unknown_api_error", "we encountered an error calling the API", map[string]interface{}{
			"response": response,
		})
		return nil, &APIError{
			HTTPCode:   http.StatusInternalServerError,
			SystemCode: "request_error",
		}

	}

	statusCode := response.StatusCode()
	if statusCode >= http.StatusMultipleChoices {
		apiError := map[string]interface{}{}
		message := "error"
		json.Unmarshal(response.Body(), &apiError)
		// check to see if there is a message field and we can throw it in the message level
		if _, ok := apiError["message"]; ok {
			message = apiError["message"].(string)
		}
		return nil, APIError{
			HTTPCode:   statusCode,
			SystemCode: "error",
			Message:    message,
			Body:       &apiError,
		}
	}

	responseData := map[string]interface{}{}
	json.Unmarshal(response.Body(), &responseData)

	return &APIReturn{
		HTTPCode: statusCode,
		Body:     responseData,
	}, nil
}
