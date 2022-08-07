package service

import "github.com/google/uuid"

const (
	ATWalletSignUp   = "/sign-up"
	ATWalletSignIn   = "/sign-in"
	ATWalletPlatform = "/user/platform"
	ATWalletStellar  = "/stellar"
	ATWalletAccount  = "/account"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    uint   `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type TokenData struct {
	ExternalID string `json:"external_id"`
	FirsName   string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}
type Asset struct {
	Balance string `json:"balance"`
	Asset   struct {
		Platform  string `json:"platform"`
		Code      string `json:"code"`
		Name      string `json:"name"`
		MinorUnit uint   `json:"minor_unit"`
	} `json:"asset"`
}
type CreateWalletResponse struct {
	Platform    string    `json:"platform"`
	Type        string    `json:"type"`
	GUID        uuid.UUID `json:"guid"`
	Name        string    `json:"name"`
	AssetsTotal uint      `json:"assets_total"`
	Assets      []Asset   `json:"assets"`
	Activated   bool      `json:"activated"`
	Registered  bool      `json:"registered"`
	Blocked     bool      `json:"blocked"`
}
type FPFPaymentResponse struct {
	Transaction string `json:"transaction"`
	Action      struct {
		IsRedirect bool   `json:"is_redirect"`
		Action     string `json:"action"`
	} `json:"action"`
}
