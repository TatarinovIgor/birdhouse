package service

import (
	"github.com/google/uuid"
)

const (
	ATWalletSignUp              = "/sign-up"
	ATWalletSignIn              = "/sign-in"
	ATWalletUserPlatform        = "/user/platform"
	ATWalletStellar             = "/stellar"
	ATWalletAccount             = "/account"
	ATWalletFPF                 = "/fpf"
	ATWalletDepositTransaction  = "/deposit/transactions"
	ATWalletWithdrawTransaction = "/withdraw/transactions"
	ATWalletTransactions        = "/transactions"
	ATWalletTomlFile            = ".well-known/stellar.toml"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    uint   `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type TokenData struct {
	Payload   UserData `json:"payload"`
	Subject   string   `json:"sub"`
	IssuedAt  uint     `json:"iat"`
	ExpiresIn uint     `json:"exp"`
}

type UserData struct {
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
	Balance       float64 `json:"balance"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	MinorUnit     uint    `json:"minor_unit"`
	Activated     bool    `json:"activated"`
	StellarCode   string  `json:"stellar_code"`
	StellarIssuer string  `json:"stellar_issuer"`
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

type DepositResponse struct {
	ID               string `json:"id"`
	StellarAccountID string `json:"stellar_account_id"`
	StellarMemoType  string `json:"stellar_memo_type"`
	StellarMemo      string `json:"stellar_memo"`
}

type DepositData struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	AmountIn         string `json:"amount_in"`
	AmountOut        string `json:"amount_out"`
	AmountFee        string `json:"amount_fee"`
	StellarAccountID string `json:"stellar_account_id"`
	StellarMemoType  string `json:"stellar_memo_type"`
	StellarMemo      string `json:"stellar_memo"`
	StartedAt        string `json:"started_at"`
	CompletedAt      string `json:"completed_at"`
}

type TransactionRequest struct {
	Transactions []TransactionData `json:"transactions"`
}
type TransactionData struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	AmountIn         string `json:"amount_in"`
	AmountOut        string `json:"amount_out"`
	AmountFee        string `json:"amount_fee"`
	StellarAccountID string `json:"stellar_account_id"`
	StellarMemoType  string `json:"stellar_memo_type"`
	StellarMemo      string `json:"stellar_memo"`
	StartedAt        string `json:"started_at"`
	Type             string `json:"type_operation"`
}

type EmailData struct {
	Email          string `json:"email"`
	EmailConfirmed bool   `json:"email_confirmed"`
}

type Authentication struct {
	AuthenticationEmail EmailData `json:"email"`
}

type UsersData struct {
	CreationDate   string         `json:"Created Date"`
	ModifiedDate   string         `json:"Modified Date"`
	UserSignedUp   bool           `json:"user_signed_up"`
	Authentication Authentication `json:"authentication"`
	FirstName      string         `json:"first_name_text"`
	LastName       string         `json:"last_name_text"`
	PhoneNumber    string         `json:"phone_text"`
	Guid           string         `json:"guid_text"`
	Balance        int            `json:"balance_number"`
	Id             string         `json:"_id"`
}

type BubbleUsersResponse struct {
	Cursor    int         `json:"cursor"`
	Result    []UsersData `json:"results"`
	Count     int         `json:"count"`
	Remaining int         `json:"remaining"`
}

type BubbleUsersData struct {
	Response BubbleUsersResponse `json:"response"`
}

type TomlConfig struct {
	Version   string `toml:"VERSION"`
	KYCServer string `toml:"KYC_SERVER"`
}
