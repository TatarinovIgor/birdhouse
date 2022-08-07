package service

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    uint   `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type TokenData struct {
	ExternalID string `json:"external_id"`
	FirsName   string `json:"firs_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}
