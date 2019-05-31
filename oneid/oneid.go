package oneid

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/inetspa/oneplatform-sdk-go/requests"
	"net/http"
)

type Api struct {
	apiEndpointBase string
	clientID        string
	clientSecret    string
	refCode         string
	callbackUrl     string
	headers         map[string]string
}

func New(clientID string, clientSecret string, refCode string, callbackUrl string) *Api {
	headers := map[string]string{
		requests.ContentType: "application/json",
		requests.UserAgent:   fmt.Sprintf(userAgent, version),
	}
	api := Api{
		apiEndpointBase: apiEndpointBase,
		clientID:        clientID,
		clientSecret:    clientSecret,
		refCode:         refCode,
		callbackUrl:     callbackUrl,
		headers:         headers,
	}
	return &api
}

func (api *Api) GetVersion() string {
	return version
}

func (api *Api) SetEndpoint(endpoint string) {
	api.apiEndpointBase = endpoint
}

func (api *Api) Login(username string, password string) (AuthenticationResult, error) {
	var result AuthenticationResult
	url := fmt.Sprintf("%s%s", api.apiEndpointBase, apiEndpointLogin)
	reqJson, err := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	}{
		ClientID:     api.clientID,
		ClientSecret: api.clientSecret,
		GrantType:    grantTypePassword,
		Username:     username,
		Password:     password,
	})
	if err != nil {
		return result, err
	}
	rawResponse, err := requests.Post(url, api.headers, bytes.NewBuffer(reqJson))
	if err := checkError(rawResponse, err); err != nil {
		return result, err
	}
	err = json.Unmarshal(rawResponse.Body, &result)
	return result, err
}

func (api *Api) GetProfile(tokenType string, accessToken string) (AccountProfile, error) {
	var profile AccountProfile
	url := fmt.Sprintf("%s%s", api.apiEndpointBase, apiEndpointGetProfile)
	if tokenType == "" || accessToken == "" {
		return profile, errors.New("login required")
	}
	headers := api.headers
	headers[requests.Authorization] = fmt.Sprintf("%s %s", tokenType, accessToken)
	rawResponse, err := requests.Get(url, headers)
	if err := checkError(rawResponse, err); err != nil {
		return profile, err
	}
	return profile, json.Unmarshal(rawResponse.Body, &profile)
}

func (api *Api) RefreshNewToken(refreshToken string) (AuthenticationResult, error) {
	var result AuthenticationResult
	url := fmt.Sprintf("%s%s", api.apiEndpointBase, apiEndpointRefreshToken)
	reqJson, err := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RefreshToken string `json:"refresh_token"`
	}{
		ClientID:     api.clientID,
		ClientSecret: api.clientSecret,
		GrantType:    grantTypeRefreshToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return result, err
	}
	rawResponse, err := requests.Post(url, api.headers, bytes.NewBuffer(reqJson))
	if err := checkError(rawResponse, err); err != nil {
		return result, err
	}
	err = json.Unmarshal(rawResponse.Body, &result)
	return result, err
}

func (api *Api) VerifyAuthorizationCode(code string) (AuthenticationResult, error) {
	var result AuthenticationResult
	if code == "" {
		return result, errors.New("authorization code required")
	}
	url := fmt.Sprintf("%s%s", api.apiEndpointBase, apiEndpointVerifyCode)
	reqJson, err := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		Scope        string `json:"scope"`
	}{
		GrantType:    grantTypeCode,
		ClientID:     api.clientID,
		ClientSecret: api.clientSecret,
		Code:         code,
		Scope:        "",
	})
	if err != nil {
		return result, err
	}
	rawResponse, err := requests.Post(url, api.headers, bytes.NewBuffer(reqJson))
	if err := checkError(rawResponse, err); err != nil {
		return result, err
	}
	err = json.Unmarshal(rawResponse.Body, &result)
	return result, err
}

func (api *Api) GetLoginUrl() string {
	loginUri := fmt.Sprintf(apiEndpointLoginUrl, api.clientID, grantTypeCode, "")
	return fmt.Sprintf("%s%s", api.apiEndpointBase, loginUri)
}

func (api *Api) RedirectToLoginUrl(w http.ResponseWriter, r *http.Request) {
	url := api.GetLoginUrl()
	http.Redirect(w, r, url, http.StatusFound)
}
