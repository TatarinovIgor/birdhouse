package service

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	stellar "github.com/stellar/go/txnbuild"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type ATWalletService struct {
	baseWalletURL    string
	requestPublicKey interface{}
	appGUID          uuid.UUID
	tokenTimeToLive  int64
}

func (service ATWalletService) SignUp(token string) (*AuthResponse, error) {
	URL := service.getATWalletUrl() + ATWalletSignUp
	// generate session and jwt from uid
	sessionID, _ := uuid.NewUUID()
	return service.authATWallet(token, sessionID.String(), URL)
}

func (service ATWalletService) CreateStellarWallet(jwtToken, token, accountType, name string) (*UserAccount, error) {
	URL := service.getATWalletUrl() + ATWalletUserPlatform + ATWalletStellar + ATWalletAccount
	body := fmt.Sprintf("{\"platfrom\": \"stellar\", \"type\": \"%s\", \"name\": \"%s\"}", accountType, name)
	session, _ := uuid.NewUUID()
	result, err := service.requestToATWallet(URL, "POST", jwtToken, token, session.String(), []byte(body))
	if err != nil {
		return nil, err
	}
	createdWallet := UserAccount{}
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

func (service ATWalletService) TokenDecode(token string) (*UserData, error) {
	tokenData := TokenData{}

	tok, err := jwt.Parse([]byte(token), jwt.WithVerify(jwa.RS256, service.requestPublicKey.(*rsa.PublicKey)), jwt.WithValidate(true))
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
	fmt.Println(time.Now().Unix(), tok.IssuedAt().Unix(), service.tokenTimeToLive, time.Now().Unix()-tok.IssuedAt().Unix())
	if time.Now().Unix()-tok.IssuedAt().Unix() > service.tokenTimeToLive {
		return nil, fmt.Errorf("token had expaired")
	}
	return &tokenData.Payload, nil
}

func (service ATWalletService) FPFPayment(jwtToken, token, assetCode, asseIssuer string,
	account uuid.UUID, amount float64, isDeposit bool) (*FPFPaymentResponse, error) {
	URL := service.getATWalletUrl() + ATWalletUserPlatform + ATWalletStellar + ATWalletAccount + "/" +
		account.String() + ATWalletFPF
	body := fmt.Sprintf("{\"amount\": %v, \"asset_code\": \"%s\", \"asset_issuer\": \"%s\", \"is_deposit\": %v}",
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

func (service ATWalletService) Deposit(jwtToken, token, assetCode, asseIssuer, accGuid, receiverExternalId, senderInternalId string,
	amount float64) (*DepositData, error) {
	URL := service.getATWalletUrl() + ATWalletUserPlatform + ATWalletStellar + ATWalletAccount + "/" +
		accGuid + ATWalletDepositTransaction
	body := fmt.Sprintf("{\"amount\": \"%v\", \"asset_code\": \"%s\", \"asset_issuer\": \"%s\", \"sender_id\": \"%s\", \"receiver_id\": \"%s\"}",
		amount, assetCode, asseIssuer, senderInternalId, receiverExternalId)
	session, _ := uuid.NewUUID()
	result, err := service.requestToATWallet(URL, "POST", jwtToken, token, session.String(), []byte(body))
	if err != nil {
		return nil, err
	}
	depositResponse := DepositResponse{}

	err = json.NewDecoder(result).Decode(&depositResponse)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var memoBytes []byte
	if memoBytes, err = base64.StdEncoding.DecodeString(depositResponse.StellarMemo); err != nil || len(memoBytes) != 32 {
		// error to log
		return nil, err
	}
	var buf [32]byte
	copy(buf[:], memoBytes)
	memo := buf
	_, err = service.BuildStellarTransactionHash("SBBYW2G6H26Y7JHUJUV2DTS7EZDE7DE5S5S4BAVH7UIST3JD26OQ4VQ6", asseIssuer, depositResponse.StellarAccountID, memo, amount)
	if err != nil {
		return nil, err
	}

	transactionData, err := service.requestToATWallet(URL+"/"+depositResponse.ID, "GET", jwtToken, token, session.String(), nil)

	depositData := DepositData{}
	err = json.NewDecoder(transactionData).Decode(&depositData)
	if err != nil {
		return nil, err
	}

	return &depositData, nil
}

func (service ATWalletService) Withdraw(jwtToken, token, assetCode, asseIssuer, accGuid, senderExternalId, receiverInternalId string,
	amount float64) (*DepositData, error) {
	URL := service.getATWalletUrl() + ATWalletUserPlatform + ATWalletStellar + ATWalletAccount + "/" +
		accGuid + ATWalletWithdrawTransaction
	session, _ := uuid.NewUUID()
	memo := generateRandom(28)
	body := fmt.Sprintf("{\"amount\": \"%v\", \"asset_code\": \"%s\", \"asset_issuer\": \"%s\", \"sender_id\": \"%s\", \"receiver_id\": \"%s\", \"memo_type\": \"%s\", \"memo\": \"%s\"}",
		amount, assetCode, asseIssuer, senderExternalId, receiverInternalId, "text", memo)
	result, err := service.requestToATWallet(URL, "POST", jwtToken, token, session.String(), []byte(body))
	if err != nil {
		return nil, err
	}
	depositResponse := DepositResponse{}

	err = json.NewDecoder(result).Decode(&depositResponse)
	if err != nil {
		return nil, err
	}

	_, err = service.BuildStellarTransactionText("SBBYW2G6H26Y7JHUJUV2DTS7EZDE7DE5S5S4BAVH7UIST3JD26OQ4VQ6", asseIssuer, depositResponse.StellarAccountID, memo, amount)
	if err != nil {
		return nil, err
	}

	transactionData, err := service.requestToATWallet(URL+"/"+depositResponse.ID, "GET", jwtToken, token, session.String(), nil)

	depositData := DepositData{}
	err = json.NewDecoder(transactionData).Decode(&depositData)
	if err != nil {
		return nil, err
	}

	return &depositData, nil
}

func (service ATWalletService) GetBalance(jwtToken, token string) (*UserPlatformResponse, error) {
	URL := service.getATWalletUrl() + ATWalletUserPlatform + ATWalletStellar
	queryParam := "?include_accounts=true&include_assets=true"
	session, _ := uuid.NewUUID()
	result, err := service.requestToATWallet(URL+queryParam, "GET", jwtToken, token,
		session.String(), nil)
	if err != nil {
		return nil, err
	}
	userPlatformResponse := UserPlatformResponse{}
	err = json.NewDecoder(result).Decode(&userPlatformResponse)
	if err != nil {
		return nil, err
	}
	return &userPlatformResponse, nil
}

func NewATWalletService(baseWalletURL string, requestPublicKey interface{}, appGUID uuid.UUID, tokenTimeToLive int64) *ATWalletService {
	return &ATWalletService{
		baseWalletURL:    baseWalletURL,
		requestPublicKey: requestPublicKey,
		appGUID:          appGUID,
		tokenTimeToLive:  tokenTimeToLive,
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

func (service ATWalletService) BuildStellarTransactionHash(addressOrSeed, assetIssuer, destination string, memo [32]byte, amount float64) (io.ReadCloser, error) {
	kp, err := keypair.Parse(addressOrSeed)
	if err != nil {
		return nil, fmt.Errorf("can't parse sender key, error: %v", err)
	}
	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return nil, fmt.Errorf("can't source sender account, error: %v", err)
	}
	op := stellar.Payment{
		Destination:   destination,
		Amount:        fmt.Sprintf("%f", amount),
		Asset:         stellar.CreditAsset{Code: "ATUSD", Issuer: assetIssuer},
		SourceAccount: sourceAccount.AccountID,
	}
	tx, err := stellar.NewTransaction(
		stellar.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []stellar.Operation{&op},
			BaseFee:              stellar.MinBaseFee,
			Preconditions:        stellar.Preconditions{TimeBounds: stellar.NewInfiniteTimeout()}, // Use a real timeout in production!
			Memo:                 stellar.MemoHash(memo),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("can't build transaction, error: %v", err)
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full))
	if err != nil {
		return nil, fmt.Errorf("can't sign transaction, error: %v", err)
	}
	txe, err := tx.Base64()
	if err != nil {
		return nil, fmt.Errorf("can't convert to base 64, error: %v", err)
	}
	fmt.Println(txe)

	resp, err := horizonclient.DefaultTestNetClient.SubmitTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("can't submit transaction, error: %v", err)
	}
	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

	return nil, nil
}

func (service ATWalletService) BuildStellarTransactionText(addressOrSeed, assetIssuer, destination, memo string, amount float64) (io.ReadCloser, error) {
	kp, err := keypair.Parse(addressOrSeed)
	if err != nil {
		return nil, fmt.Errorf("can't parse sender key, error: %v", err)
	}
	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return nil, fmt.Errorf("can't source sender account, error: %v", err)
	}
	op := stellar.Payment{
		Destination:   destination,
		Amount:        fmt.Sprintf("%f", amount),
		Asset:         stellar.CreditAsset{Code: "ATUSD", Issuer: assetIssuer},
		SourceAccount: sourceAccount.AccountID,
	}
	tx, err := stellar.NewTransaction(
		stellar.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []stellar.Operation{&op},
			BaseFee:              stellar.MinBaseFee,
			Preconditions:        stellar.Preconditions{TimeBounds: stellar.NewInfiniteTimeout()}, // Use a real timeout in production!
			Memo:                 stellar.MemoText(memo),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("can't build transaction, error: %v", err)
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full))
	if err != nil {
		return nil, fmt.Errorf("can't sign transaction, error: %v", err)
	}
	txe, err := tx.Base64()
	if err != nil {
		return nil, fmt.Errorf("can't convert to base 64, error: %v", err)
	}
	fmt.Println(txe)

	resp, err := horizonclient.DefaultTestNetClient.SubmitTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("can't submit transaction, error: %v", err)
	}
	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

	return nil, nil
}

func (service ATWalletService) convertMemo(memoString string) ([]byte, error) {
	txHashBytes, err := base64.StdEncoding.DecodeString(memoString)
	if err != nil {
		return nil, fmt.Errorf("can't convert memo from base64, err: %s", err)
	}
	return txHashBytes, nil
}

func generateRandom(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
