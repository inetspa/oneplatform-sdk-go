# OnePlatform API SDK for Go

![OnePlatform Mascot](https://monitor.sdi.one.th/imagik/bj612eatpnstbpk7nhsg)

## เกี่ยวกับ OnePlatform API

ติดตามคู่มือการใช้งาน Official API และข้อมูลอื่นๆ ได้ที่ [https://api.one.th]

## Requirements

This library requires Go 1.10 or later.

## การติดตั้ง

```sh
$ go get github.com/inetspa/oneplatform-sdk-go/...
```

## ตัวอย่างการใช้งาน

```go
import (
    "github.com/inetspa/oneplatform-sdk-go/oneid"
)

func main() {
    id := oneid.New("_CLIENT_ID_", "_CLIENT_SECRET", "_REF_CODE_", "_CALLBACK_URL_")
    r, err := id.Login("_USERNAME_", "_PASSWORD_")
    ...
}
```

## Platform services list
* One ID:
    * Login with username and password
    * Get profile
    * Verify authorization code
    * Refresh token
    * Generate login link
    * Redirect to login link

## APIs

### One ID

#### Login (OAuth2 - Password Grant)
```go
r, err := id.Login("_USERNAME_", "_PASSWORD_")
if err != nil {
    // Do something when login failed.
}
```

### Refresh Token
```go
c, err := id.RefreshNewToken(r.RefreshToken)
if err != nil {
    // Do something when refresh failed.
}
```

### Verify authorization code (OAuth2 - Authorization code)
```go
r, err := id.VerifyAuthorizationCode("_AUTHORIZATION_CODE_")
if err != nil {
    // Do something when verify code failed.
}
```

#### Authentication result structure
- Login
- Refresh token
- Verify authorization code
```go
type AuthenticationResult struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccountID    string `json:"account_id"`
	Result       string `json:"result"`
	Username     string `json:"username"`
}
```

### Get account profile
```go
p, err := id.GetProfile(r.TokenType, r.AccessToken)
if err != nil {
    // Do something when cannot get account profile.
}
```

### Generate login link
```go
url := id.GetLoginUrl()
```

### Redirect to login url
```go
id.RedirectToLoginUrl(w, r)
```

## Changelog

### Version 0.1.2 (2019-05-31)

* Bug fixed
    * Fixed verify authorization code error

### Version 0.1.0 (2019-05-07)

* First release after prepare library
* Platform service compatibility:
    * One ID:
        * Login with username and password
        * Get profile
        * Verify authorization code
        * Refresh token
        * Generate login link
        * Redirect to login link

[https://api.one.th]: <https://api.one.th>
