package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
)

type ATWalletService struct {
	baseWalletURL    string
	requestPublicKey string
	appGUID          uuid.UUID
}

func (service ATWalletService) SignUp(token string) (*AuthResponse, error) {
	URL := service.getATWalletUrl() + "/sign-up"
	// generate session and jwt from uid
	sessionID, _ := uuid.NewUUID()
	return service.authATWallet(token, sessionID.String(), URL)
}

func (service ATWalletService) Activate(token string) (*AuthResponse, error) {
	return nil, nil
}

func (service ATWalletService) SignIn(token string) (*AuthResponse, error) {
	URL := service.getATWalletUrl() + "/sign-in"
	// generate session and jwt from uid
	sessionID, _ := uuid.NewUUID()
	return service.authATWallet(token, sessionID.String(), URL)
}

func (service ATWalletService) TokenDecode(token string) (*TokenData, error) {
	tokenData := TokenData{}
	tok, err := jwt.Parse([]byte(token), jwt.WithVerify(false), jwt.WithValidate(false))
	if err != nil {
		return nil, fmt.Errorf("can't parse token, err: %v", err)
	}
	tokenString, err := json.Marshal(tok.PrivateClaims())
	if err != nil {
		return nil, fmt.Errorf("can't marshal token claim, err: %v", err)
	}
	err = json.Unmarshal(tokenString, &tokenData)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshak token data, err: %v", err)
	}
	return &tokenData, nil
}

func NewATWalletService(baseWalletURL, requestPublicKey string, appGUID uuid.UUID) *ATWalletService {
	return &ATWalletService{
		baseWalletURL:    baseWalletURL,
		requestPublicKey: requestPublicKey,
		appGUID:          appGUID,
	}
}

func (service ATWalletService) getATWalletUrl() string {
	return service.baseWalletURL + "/" + service.appGUID.String()
}
func (service ATWalletService) authATWallet(token, session, url string) (*AuthResponse, error) {
	authResponse := AuthResponse{}
	// create headers
	client := http.Client{}
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("can't make request for signup, err %v", err)
	}
	request.Header.Set("X-Auth-Token", token)
	request.Header.Set("X-Session-ID", session)
	// make a request
	result, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("can't request for signup, err %v", err)
	}
	if result.StatusCode/100 != 2 {
		return nil, fmt.Errorf("unexpected code from auth url: %s, code: %v", url, result.StatusCode)
	}
	err = json.NewDecoder(result.Body).Decode(&authResponse)
	if err != nil {
		return nil, fmt.Errorf("can't parse token from wallet %v", err)
	}
	return &authResponse, nil
}
