package hubspot

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty"
)

var oAuthToken = OAuthToken{}

// OAuthToken represents a current oAuth token for accessing calls requiring oAuth
type OAuthToken struct {
	RefreshToken string
	AccessToken  string
	ExpiresIn    int
}

// OAuthResponse represents the oAuth response when trying to refresh a token
type OAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func renewOAuthToken() (*OAuthToken, error) {
	resty.SetContentLength(true)
	response, err := resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": Config.HubSpotOAuthRefreshToken,
			"client_id":     Config.HubspotClientID,
			"client_secret": Config.HubspotClientSecret,
		}).
		Post("https://api.hubapi.com/oauth/v1/token")

	if err != nil {
		return nil, err
	}
	code := response.StatusCode()
	if code != 200 {
		return nil, errors.New("Could not refresh that token")
	}
	parsedResponse := OAuthResponse{}
	json.Unmarshal(response.Body(), &parsedResponse)
	newToken := OAuthToken{
		AccessToken:  parsedResponse.AccessToken,
		RefreshToken: parsedResponse.RefreshToken,
		ExpiresIn:    parsedResponse.ExpiresIn,
	}
	oAuthToken = newToken
	return &newToken, nil
}

func init() {
	// get a new oauth token
	_, err := renewOAuthToken()
	if err != nil {
		log("error", "hubspot_sdk_no_oauth", "could not get oauth token; you will not be able to use any calls that require oAuth", err)
	}
}
