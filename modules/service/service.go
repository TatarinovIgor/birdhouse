package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"io"
	"net/http"
)

type ATWalletService struct {
	baseWalletURL    string
	requestPublicKey string
	appGUID          uuid.UUID
}

func (service ATWalletService) SignUp(token string) (*AuthResponse, error) {
	URL := service.getATWalletUrl() + ATWalletSignUp
	// generate session and jwt from uid
	sessionID, _ := uuid.NewUUID()
	return service.authATWallet(token, sessionID.String(), URL)
}

func (service ATWalletService) CreateStellarWallet(jwtToken, token, accountType, name string) (*CreateWalletResponse, error) {
	URL := service.getATWalletUrl() + ATWalletPlatform + ATWalletStellar + ATWalletAccount
	body := fmt.Sprintf("{\"platfrom\": \"stellar\", \"type\": \"%s\", \"name\": \"%s\"}", accountType, name)
	session, _ := uuid.NewUUID()
	result, err := service.requestToATWallet(URL, "POST", jwtToken, token, session.String(), []byte(body))
	if err != nil {
		return nil, err
	}
	createdWallet := CreateWalletResponse{}
	err = json.NewDecoder(result).Decode(&createdWallet)
	if err != nil {
		return nil, err
	}
	return &createdWallet, nil
}

func (service ATWalletService) SignIn(token string) (*AuthResponse, error) {
	URL := service.getATWalletUrl() + ATWalletSignIn
	// generate session and jwt from uid
	sessionID, _ := uuid.NewUUID()
	return service.authATWallet(token, sessionID.String(), URL)
}

func (service ATWalletService) TokenDecode(token string) (*TokenData, error) {
	tokenData := TokenData{}
	tok, err := jwt.Parse([]byte(token), jwt.WithVerify(false), jwt.WithValidate(false))
	if err != nil {
		return nil, fmt.Errorf("can't parse token: %s, err: %v", token, err)
	}
	tokenString, err := json.Marshal(tok.PrivateClaims())
	if err != nil {
		return nil, fmt.Errorf("can't marshal token claim, err: %v", err)
	}
	err = json.Unmarshal(tokenString, &tokenData)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal token data, err: %v", err)
	}
	return &tokenData, nil
}
func (service ATWalletService) FPFPayment(jwtToken, token, assetCode, asseIssuer string,
	account uuid.UUID, amount float64, isDeposit bool) (*FPFPaymentResponse, error) {
	URL := service.getATWalletUrl() + ATWalletPlatform + ATWalletStellar + ATWalletAccount
	body := fmt.Sprintf("{\"amount\": %v, \"asset_code\": \"%s\", \"asset_issuer\": \"%s\", \"is_depoist\": %v}",
		amount, assetCode, asseIssuer, isDeposit)
	session, _ := uuid.NewUUID()
	result, err := service.requestToATWallet(URL, "POST", jwtToken, token, session.String(), []byte(body))
	if err != nil {
		return nil, err
	}
	paymentResponse := FPFPaymentResponse{}
	err = json.NewDecoder(result).Decode(&paymentResponse)
	if err != nil {
		return nil, err
	}
	return &paymentResponse, nil
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
		return nil, fmt.Errorf("can't make request for auth url: %s, err %v", url, err)
	}
	request.Header.Set("X-Auth-Token", token)
	request.Header.Set("X-Session-ID", session)
	// make a request
	result, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("can't request for auth url: %s, err %v", url, err)
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
func (service ATWalletService) requestToATWallet(url, requestType, jwtToken, token, session string,
	body []byte) (io.ReadCloser, error) {
	// create headers
	client := http.Client{}
	request, err := http.NewRequest(requestType, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("can't make %s request for url: %s, err %v", requestType, url, err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-Auth-Token", jwtToken)
	request.Header.Set("X-Session-ID", session)

	// make a request
	result, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("can't request %s for url: %s, err %v", requestType, url, err)
	}
	if result.StatusCode/100 != 2 {
		defer func() { _ = result.Body.Close() }()
		bodyBytes, _ := io.ReadAll(result.Body)
		return nil, fmt.Errorf("unexpected code from %s request for url: %s, code %v, message: %s",
			requestType, url, result.StatusCode, string(bodyBytes))
	}
	return result.Body, nil
}
