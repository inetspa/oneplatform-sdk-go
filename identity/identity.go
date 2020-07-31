package identity

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/inetspa/golib/requests"
	"github.com/inetspa/golib/web"
	"io"
	"net/http"
)

const (
	// Api
	apiEndpoint = "https://one.th"

	// Information
	grantTypePassword     = "password"
	grantTypeCode         = "authorization_code"
	grantTypeRefreshToken = "refresh_token"
)

func NewIdentity(clientID string, clientSecret string, refCode string, callbackUrl string) *Identity {
	headers := map[string]string{
		web.HeaderContentType: web.MIMEApplicationJSON,
	}
	id := Identity{
		apiEndpoint:  apiEndpoint,
		clientId:     clientID,
		clientSecret: clientSecret,
		refCode:      refCode,
		callbackUrl:  callbackUrl,
		headers:      headers,
	}
	return &id
}

func (id *Identity) Login(username string, password string) (AuthenticationResult, error) {
	var result AuthenticationResult
	reqJson, err := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	}{
		ClientID:     id.clientId,
		ClientSecret: id.clientSecret,
		GrantType:    grantTypePassword,
		Username:     username,
		Password:     password,
	})
	if err != nil {
		return result, err
	}
	r, err := requests.Post(id.url("/api/oauth/getpwd"), id.headers, bytes.NewBuffer(reqJson), requests.Timeout)
	if err != nil {
		return result, err
	}
	if r.Code != http.StatusOK {
		return result, errors.New(fmt.Sprintf("client return error with code %d (%s)", r.Code, string(r.Body)))
	}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (id *Identity) GetProfile(tokenType string, accessToken string) (AccountProfile, error) {
	var profile AccountProfile
	if tokenType == "" || accessToken == "" {
		return profile, errors.New("login required")
	}
	headers := id.headers
	headers[web.HeaderAuthorization] = fmt.Sprintf("%s %s", tokenType, accessToken)
	r, err := id.send(http.MethodGet, id.url("/api/account"), nil, headers)
	if err != nil {
		return profile, err
	}
	return profile, json.Unmarshal(r.Body, &profile)
}

func (id *Identity) RefreshNewToken(refreshToken string) (AuthenticationResult, error) {
	var result AuthenticationResult
	reqJson, err := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RefreshToken string `json:"refresh_token"`
	}{
		ClientID:     id.clientId,
		ClientSecret: id.clientSecret,
		GrantType:    grantTypeRefreshToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return result, err
	}
	r, err := id.send(http.MethodPost, id.url("/api/oauth/get_refresh_token"), bytes.NewBuffer(reqJson), nil)
	if err != nil {
		return result, err
	}
	if r.Code != http.StatusOK {
		return result, errors.New(fmt.Sprintf("client return error with code %d (%s)", r.Code, string(r.Body)))
	}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (id *Identity) VerifyAuthorizationCode(code string) (AuthenticationResult, error) {
	var result AuthenticationResult
	if code == "" {
		return result, errors.New("authorization code required")
	}
	reqJson, err := json.Marshal(&struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		Scope        string `json:"scope"`
	}{
		GrantType:    grantTypeCode,
		ClientID:     id.clientId,
		ClientSecret: id.clientSecret,
		Code:         code,
		Scope:        "",
	})
	if err != nil {
		return result, err
	}
	r, err := id.send(http.MethodPost, id.url("/oauth/token"), bytes.NewBuffer(reqJson), nil)
	if err != nil {
		return result, err
	}
	if r.Code != http.StatusOK {
		return result, errors.New(fmt.Sprintf("client return error with code %d (%s)", r.Code, string(r.Body)))
	}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (id *Identity) GetLoginUrl() string {
	return id.url(fmt.Sprintf("/api/oauth/getcode?client_id=%s&response_type=%s&scope=%s", id.clientId, grantTypeCode, ""))
}

func (id *Identity) RedirectToLoginUrl(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, id.GetLoginUrl(), http.StatusFound)
}

func (id *Identity) SetEndpoint(endpoint string) {
	id.apiEndpoint = endpoint
}

func (id *Identity) url(path string) string {
	return fmt.Sprintf("%s%s", id.apiEndpoint, path)
}

func (id *Identity) send(method string, url string, body io.Reader, headers map[string]string) (requests.Response, error) {
	if headers != nil {
		return requests.Request(method, url, headers, body, requests.Timeout)
	}
	return requests.Request(method, url, id.headers, body, requests.Timeout)
}
