package oneid

const (
	// Base
	apiEndpointBase = "https://one.th"

	// Authentication
	apiEndpointLogin        = "/api/oauth/getpwd"
	apiEndpointRefreshToken = "/api/oauth/get_refresh_token"
	apiEndpointVerifyCode   = "/oauth/token"
	apiEndpointLoginUrl     = "/api/oauth/getcode?client_id=%s&response_type=%s&scope=%s"

	// Profile
	apiEndpointGetProfile = "/api/account"

	// Information
	userAgent             = "oneplatform-sdk-go/%s"
	grantTypePassword     = "password"
	grantTypeCode         = "authorization_code"
	grantTypeRefreshToken = "refresh_token"
)
