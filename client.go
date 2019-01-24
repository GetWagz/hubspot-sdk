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

func makeCall(method, endpoint string, data interface{}) (ret *APIReturn, err error) {
	if strings.HasPrefix(endpoint, "/") {
		endpoint = endpoint[1:]
	}

	// TODO: mocking

	url := fmt.Sprintf("%s%s", Config.RootURL, endpoint)
	fmt.Printf("\n===================================\n%s\n", url)

	var response *resty.Response

	request := resty.R().
		SetHeader("Accept", "application/json")

	queryParams := map[string]string{
		"hapikey": Config.HubspotAPIKey,
	}

	// Now, do what we need to do depending on the method
	var reqErr error

	switch method {
	case http.MethodGet:
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
		response, reqErr = request.SetQueryParams(queryParams).Get(url)
	case http.MethodDelete:
		response, reqErr = request.SetQueryParams(queryParams).Delete(url)
	case http.MethodPost:
		response, reqErr = request.SetQueryParams(queryParams).SetBody(data).Post(url)
	}

	if reqErr != nil {
		return nil, &APIError{
			HTTPCode:   http.StatusInternalServerError,
			SystemCode: "request_error",
			Message:    reqErr.Error(),
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
