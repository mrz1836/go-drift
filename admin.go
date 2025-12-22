package drift

// TokenInfoRequest represents the request body for token info
type TokenInfoRequest struct {
	AccessToken string `json:"access_token"`
}

// TokenInfo represents the token information response
type TokenInfo struct {
	AccessToken         string `json:"access_token"`
	AuthenticatedUserID string `json:"authenticated_userid"`
	CreatedAt           int64  `json:"created_at"`
	CredentialID        string `json:"credential_id"`
	ExpiresIn           int64  `json:"expires_in"`
	ID                  string `json:"id"`
	Scope               string `json:"scope"`
	TokenType           string `json:"token_type"`
}
