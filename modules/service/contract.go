package service

import "github.com/google/uuid"

const (
	ATWalletSignUp       = "/sign-up"
	ATWalletSignIn       = "/sign-in"
	ATWalletUserPlatform = "/user/platform"
	ATWalletStellar      = "/stellar"
	ATWalletAccount      = "/account"
	ATWalletFPF          = "/fpf"
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
	Platform  string `json:"platform"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	MinorUnit uint   `json:"minor_unit"`
}
type PlatformAccount struct {
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
type UserAsset struct {
	Balance string  `json:"balance"`
	Assets  []Asset `json:"assets"`
}
type UserAccount struct {
	Platform    string      `json:"platform"`
	Type        string      `json:"type"`
	GUID        uuid.UUID   `json:"guid"`
	Name        string      `json:"name"`
	AssetsTotal uint        `json:"assets_total"`
	Assets      []UserAsset `json:"assets"`
	Activated   bool        `json:"activated"`
	Registered  bool        `json:"registered"`
	Blocked     bool        `json:"blocked"`
}
type UserPlatformResponse struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	AccountsTotal uint          `json:"accounts_total"`
	Account       []UserAccount `json:"accounts"`
}
type PlatformResponse struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	AccountsTotal uint              `json:"accounts_total"`
	Account       []PlatformAccount `json:"accounts"`
}

type FPFPaymentResponse struct {
	Transaction string `json:"transaction"`
	Action      struct {
		IsRedirect bool   `json:"is_redirect"`
		Action     string `json:"action"`
	} `json:"action"`
}
