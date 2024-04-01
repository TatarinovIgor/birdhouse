package service

import (
	"birdhouse/modules/storage"
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	encoder "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	stellar "github.com/stellar/go/txnbuild"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ATWalletService struct {
	baseWalletURL    string
	requestPublicKey interface{}
	privateKey       interface{}
	appGUID          uuid.UUID
	tokenTimeToLive  int64
	seed             string
	store            *storage.KeysPSQL
	systemWallet     string
	systemWalletSeed string

	tokenCode              string
	tokenAsset             string
	tokenBlockchain        string
	tokenIssuer            string
	processorUrl           string
	processorAdminMerchant string
	processorAdminExternal string
}

type TokenPayload struct {
	ExternalId string `json:"external_id"`
	MerchantID string `json:"merchant_id"`
}

type MintData struct {
	Code              string `json:"code"`
	Blockchain        string `json:"blockchain"`
	Amount            string `json:"amount"`
	Issuer            string `json:"issuer"`
	Type              string `json:"type"`
	ReceivingWalletID string `json:"receiving_wallet_id"`
}

type MintTokenResponse struct {
	Issuer string
	TxHash string
}

type ProcessingCredentialDeposit struct {
	Amount     float64 `json:"amount"`
	Blockchain string  `json:"blockchain"`
	Asset      string  `json:"asset"`
	Issuer     string  `json:"issuer"`
}

type ProcessorDepositResponse struct {
	WalletAddress string `json:"wallet_address"`
	Memo          string `json:"memo"`
	URL           string `json:"url"`
	Id            string `json:"id"`
}

type ProcessingCredentialWithdraw struct {
	Amount        float64 `json:"amount"`
	Blockchain    string  `json:"blockchain"`
	WalletAddress string  `json:"wallet_address"`
	Asset         string  `json:"asset"`
	Issuer        string  `json:"issuer"`
	Memo          string  `json:"memo"`
}

type ProcessorWithdrawResponse struct {
	TransactionHash string `json:"result"`
}

type DeleteWithdrawResponse struct {
	Status string `json:"status"`
}

func (service ATWalletService) CreateToken(externalID, merchantID string) (string, error) {
	token := encoder.New(encoder.SigningMethodRS256)
	var tokenToSend TokenPayload
	tokenToSend.MerchantID = merchantID
	tokenToSend.ExternalId = externalID
	claims := token.Claims.(encoder.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["payload"] = tokenToSend
	claims["iat"] = time.Now().Unix()
	tokenString, err := token.SignedString(service.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (service ATWalletService) MintTokensOnProcessing(externalID, merchantID, amount string) error {
	token, err := service.CreateToken(externalID, merchantID)
	if err != nil {
		return err
	}
	var mintData = MintData{Code: service.tokenCode, Blockchain: service.tokenBlockchain, Amount: amount, Issuer: service.tokenIssuer, Type: "", ReceivingWalletID: externalID}
	mintDataParsed, err := json.Marshal(mintData)
	r, err := http.NewRequest("POST", service.processorUrl+"/token_mint", bytes.NewBuffer(mintDataParsed))
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", token)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	var responseMint MintTokenResponse
	err = json.NewDecoder(res.Body).Decode(&responseMint)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return err
	}
	log.Println(responseMint.TxHash, responseMint.Issuer)
	return nil
}

func (service ATWalletService) DepositTokensOnProcessing(externalID, merchantID string, amount float64) (string, error) {
	token, err := service.CreateToken(externalID, merchantID)
	if err != nil {
		return "", err
	}
	var mintData = ProcessingCredentialDeposit{Asset: service.tokenAsset, Blockchain: service.tokenBlockchain, Amount: amount, Issuer: service.tokenIssuer}
	mintDataParsed, err := json.Marshal(mintData)
	r, err := http.NewRequest("POST", service.processorUrl+"/deposit", bytes.NewBuffer(mintDataParsed))
	if err != nil {
		return "", err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", token)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	var depositResponse ProcessorDepositResponse
	err = json.NewDecoder(res.Body).Decode(&depositResponse)
	if err != nil {
		return "", err
	}
	log.Println(depositResponse.WalletAddress)
	return depositResponse.WalletAddress, nil
}

func (service ATWalletService) WithdrawTokensOnProcessing(wallet, externalID, merchantID string, amount float64) (string, error) {
	token, err := service.CreateToken(service.processorAdminExternal, service.processorAdminMerchant)
	if err != nil {
		return "", err
	}
	var mintData = ProcessingCredentialWithdraw{Asset: service.tokenAsset, Blockchain: service.tokenBlockchain, Amount: amount, Issuer: service.tokenIssuer, WalletAddress: wallet}
	mintDataParsed, err := json.Marshal(mintData)
	r, err := http.NewRequest("POST", service.processorUrl+"/withdraw", bytes.NewBuffer(mintDataParsed))
	if err != nil {
		return "", err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", token)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	var withdrawResponse ProcessorWithdrawResponse
	err = json.NewDecoder(res.Body).Decode(&withdrawResponse)
	if err != nil {
		return "", err
	}
	log.Println(withdrawResponse.TransactionHash)
	return withdrawResponse.TransactionHash, nil
}

func (service ATWalletService) UpdateWithdrawTokensOnProcessing(guid, externalID, merchantID string, amount float64) (string, error) {
	token, err := service.CreateToken(service.processorAdminExternal, service.processorAdminMerchant)
	if err != nil {
		return "", err
	}
	r, err := http.NewRequest("PUT", service.processorUrl+"/withdraw/"+guid, nil)
	if err != nil {
		return "", err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", token)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	var withdrawResponse DeleteWithdrawResponse
	err = json.NewDecoder(res.Body).Decode(&withdrawResponse)
	if err != nil {
		return "", err
	}
	log.Println(withdrawResponse.Status)
	return withdrawResponse.Status, nil
}

func (service ATWalletService) ProcessDeposit(Seed, Address, Issuer, Code, externalID, merchantID string, amount float64) error {
	stellarTransactionText, err := service.BuildStellarTransactionText(Seed, Issuer, Code, service.systemWallet, "", amount)
	if err != nil {
		return err
	}
	log.Println(stellarTransactionText)
	strAmount := strconv.Itoa(int(amount))

	err = service.MintTokensOnProcessing(externalID, merchantID, strAmount)
	if err != nil {
		//Send back funds if failed to receive on processing
		_, err := service.BuildStellarTransactionText(service.systemWalletSeed, Issuer, Code, Address, "", amount)
		if err != nil {
			return err
		}
		return err
	}

	// deposit
	wallet, err := service.DepositTokensOnProcessing(externalID, merchantID, amount)
	if err != nil {
		//Send back funds if failed to receive on processing
		_, err := service.BuildStellarTransactionText(service.systemWalletSeed, Issuer, Code, Address, "", amount)
		if err != nil {
			return err
		}
		return err
	}

	// Withdraw
	guid, err := service.WithdrawTokensOnProcessing(wallet, externalID, merchantID, amount)
	if err != nil {
		//Send back funds if failed to receive on processing
		_, err := service.BuildStellarTransactionText(service.systemWalletSeed, Issuer, Code, Address, "", amount)
		if err != nil {
			return err
		}
		return err
	}

	// update withdraw
	status, err := service.UpdateWithdrawTokensOnProcessing(guid, externalID, merchantID, amount)
	if err != nil {
		//Send back funds if failed to receive on processing
		_, err := service.BuildStellarTransactionText(service.systemWalletSeed, Issuer, Code, Address, "", amount)
		if err != nil {
			return err
		}
		return err
	}
	fmt.Println(status)
	return nil
}
func (service ATWalletService) ListenAndServe() {
	ticker := time.NewTicker(time.Second * 1)
	next := int64(0)
	for {
		select {
		case <-ticker.C:
			records, err := service.store.GetNext(next, 1)
			if err != nil || len(records) == 0 {
				next = 0
				if err != nil {
					log.Println("Error while getting DB records:", err)
				}
			} else {
				record := records[len(records)-1]
				accountRequest := horizonclient.AccountRequest{AccountID: record.Key}
				client := horizonclient.DefaultTestNetClient

				account, err := client.AccountDetail(accountRequest)
				if err != nil {
					return
				}
				for i := 0; i < len(account.Balances); i++ {
					if account.Balances[i].Code == "ATUSD" {
						log.Println("Account has:", account.Balances[i].Balance)
						wallet := Wallet{}
						err = json.Unmarshal(record.Data, &wallet)
						if err != nil {
							log.Println("Unable to parse wallet for processing, err:", err)
						}
						parsedBalance, err := strconv.ParseFloat(account.Balances[i].Balance, 64)
						if err != nil {
							log.Println("Unable to parse wallet for processing, err:", err)
						}
						if parsedBalance != 0 {
							err := service.ProcessDeposit(wallet.WalletSeed, wallet.WalletAddress, account.Balances[i].Issuer, account.Balances[i].Code, record.ExternalID, record.MerchantID, parsedBalance)
							if err != nil {
								log.Println("Unable to process deposit, err:", err)
							}
						}
					}
				}
				next = record.ID
			}
		}
	}
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
	fmt.Println(err)
	fmt.Println("succes on first err")
	fmt.Println(body)
	if err != nil {
		return nil, err
	}
	createdWallet := UserAccount{}
	err = json.NewDecoder(result).Decode(&createdWallet)
	if err != nil {
		return nil, err
	}
	fmt.Println("succes on second err")

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

func (service ATWalletService) GetTransactions(jwtToken, token string,
	account uuid.UUID) (*TransactionRequest, error) {
	URL := service.getATWalletUrl() + ATWalletUserPlatform + ATWalletStellar + ATWalletAccount + "/" +
		account.String() + ATWalletTransactions
	session, _ := uuid.NewUUID()
	result, err := service.requestToATWallet(URL, "GET", jwtToken, token, session.String(), nil)
	if err != nil {
		return nil, err
	}
	transactionResponse := TransactionRequest{}
	err = json.NewDecoder(result).Decode(&transactionResponse)
	if err != nil {
		return nil, err
	}
	return &transactionResponse, nil
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
	_, err = service.BuildStellarTransactionHash(service.seed, asseIssuer, depositResponse.StellarAccountID, memo, amount)

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
	fmt.Println(senderExternalId, receiverInternalId, jwtToken)
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

	_, err = service.BuildStellarTransactionText(service.seed, asseIssuer, assetCode, depositResponse.StellarAccountID, memo, amount)
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

func (service ATWalletService) KYCGetServer(jwtToken, token, accGuid string) (string, error) {
	URL := service.getATWalletUrl() + ATWalletUserPlatform + ATWalletStellar + ATWalletAccount + "/" +
		accGuid + ATWalletTomlFile
	session, _ := uuid.NewUUID()
	result, err := service.requestToATWallet(URL, "GET", jwtToken, token, session.String(), nil)
	if err != nil {
		return "", err
	}
	tomlBytes, err := io.ReadAll(result)
	if err != nil {
		return "", err
	}
	tomlConfig := TomlConfig{}
	_, err = toml.Decode(string(tomlBytes), &tomlConfig)
	if err != nil {
		return "", err
	}
	return tomlConfig.KYCServer, nil
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

func NewATWalletService(baseWalletURL, seed, systemWallet, systemWalletSeed, tokenCode, tokenAsset, tokenBlockchain, tokenIssuer, processorUrl, processorAdminMerchant, processorAdminExternal string, requestPublicKey interface{}, privateKey interface{}, appGUID uuid.UUID, tokenTimeToLive int64, store *storage.KeysPSQL) *ATWalletService {

	return &ATWalletService{
		baseWalletURL:          baseWalletURL,
		requestPublicKey:       requestPublicKey,
		appGUID:                appGUID,
		tokenTimeToLive:        tokenTimeToLive,
		seed:                   seed,
		store:                  store,
		systemWallet:           systemWallet,
		systemWalletSeed:       systemWalletSeed,
		tokenCode:              tokenCode,
		tokenAsset:             tokenAsset,
		tokenBlockchain:        tokenBlockchain,
		tokenIssuer:            tokenIssuer,
		privateKey:             privateKey,
		processorUrl:           processorUrl,
		processorAdminMerchant: processorAdminMerchant,
		processorAdminExternal: processorAdminExternal,
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
	response, _ := io.ReadAll(result.Body)
	if result.StatusCode/100 != 2 {
		return nil, fmt.Errorf("unexpected code from auth url: %s, code: %v", url, result.StatusCode, string(response))
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

func (service ATWalletService) BuildStellarTransactionText(addressOrSeed, assetIssuer, Code, destination, memo string, amount float64) (io.ReadCloser, error) {
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
		Asset:         stellar.CreditAsset{Code: Code, Issuer: assetIssuer},
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

type Transaction struct {
	NetworkPassphrase string `json:"network_passphrase"`
	Transaction       string `json:"transaction"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type InteractiveResp struct {
	ID  string `json:"id"`
	URL string `json:"URL"`
}

// ToDo refactor
func (service ATWalletService) CreateStellarDeposit(externalId, merchantId, blockchain string) (InteractiveResp, error) {
	if (externalId == "") || (merchantId == "") {
		return InteractiveResp{}, nil
	}

	walletAddress, walletSeed, err := service.GetExistingWallet(externalId, merchantId, blockchain)

	if err != nil {
		return InteractiveResp{}, err
	}

	wellKnownResp, err := http.Get("https://gollum-sep24.armenotech.net/gollum/api/v1/sep0024/info")
	homeDomain := "gollum-sep24.armenotech.net"
	clientDomain := "demo-wallet-server.stellar.org"
	singingKey := "GBU5DCUJV5CNGNWHG4ABGDI37AMXEQFKCRBPY7YYRYRVKTIPF35P6CJE" //ToDo fetch from toml
	if err != nil {
		return InteractiveResp{}, err
	}
	_, err = ioutil.ReadAll(wellKnownResp.Body)
	if err != nil {
		return InteractiveResp{}, err
	}
	sep0010Resp, err := http.Get("https://" + homeDomain + "/gollum/api/v1/sep0010?account=" + walletAddress + "&home_domain=" + homeDomain + "&client_domain=" + clientDomain)
	if err != nil {
		return InteractiveResp{}, err
	}
	sep0010, err := ioutil.ReadAll(sep0010Resp.Body)
	if err != nil {
		return InteractiveResp{}, err
	}
	var transaction Transaction
	err = json.Unmarshal(sep0010, &transaction)
	if err != nil {
		return InteractiveResp{}, err
	}
	token, err := SignTransaction(walletSeed, singingKey, transaction.Transaction)
	if err != nil {
		return InteractiveResp{}, err
	}

	client := &http.Client{}

	form := url.Values{}
	form.Add("asset_code", "ATUSD")
	form.Add("account", walletAddress)
	form.Add("lang", "en")
	form.Add("claimable_balance_supported", "false")

	bearer := "Bearer " + token

	interactiveReq, err := http.NewRequest("POST", "https://gollum-sep24.armenotech.net/gollum/api/v1/sep0024/transactions/deposit/interactive", strings.NewReader(form.Encode()))
	if err != nil {
		return InteractiveResp{}, err
	}
	interactiveReq.Header.Add("Authorization", bearer)
	interactiveReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(interactiveReq)
	if err != nil {
		return InteractiveResp{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	res := InteractiveResp{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return InteractiveResp{}, err
	}
	return res, nil
}

func (service ATWalletService) CreateStellarWithdraw(externalId, merchantId, blockchain string) (InteractiveResp, error) {

	walletAddress, walletSeed, err := service.GetExistingWallet(externalId, merchantId, blockchain)

	if err != nil {
		return InteractiveResp{}, err
	}

	wellKnownResp, err := http.Get("https://gollum-sep24.armenotech.net/gollum/api/v1/sep0024/info")
	homeDomain := "gollum-sep24.armenotech.net"
	clientDomain := "demo-wallet-server.stellar.org"
	singingKey := "GBU5DCUJV5CNGNWHG4ABGDI37AMXEQFKCRBPY7YYRYRVKTIPF35P6CJE" //ToDo fetch from toml
	if err != nil {
		return InteractiveResp{}, err
	}
	_, err = ioutil.ReadAll(wellKnownResp.Body)
	if err != nil {
		return InteractiveResp{}, err
	}
	sep0010Resp, err := http.Get("https://" + homeDomain + "/gollum/api/v1/sep0010?account=" + walletAddress + "&home_domain=" + homeDomain + "&client_domain=" + clientDomain)
	if err != nil {
		return InteractiveResp{}, err
	}
	sep0010, err := ioutil.ReadAll(sep0010Resp.Body)
	if err != nil {
		return InteractiveResp{}, err
	}
	var transaction Transaction
	err = json.Unmarshal(sep0010, &transaction)
	if err != nil {
		return InteractiveResp{}, err
	}
	token, err := SignTransaction(walletSeed, singingKey, transaction.Transaction)
	if err != nil {
		return InteractiveResp{}, err
	}

	client := &http.Client{}

	form := url.Values{}
	form.Add("asset_code", "ATUSD")
	form.Add("account", walletAddress)
	form.Add("lang", "en")
	form.Add("claimable_balance_supported", "false")

	bearer := "Bearer " + token

	interactiveReq, err := http.NewRequest("POST", "https://gollum-sep24.armenotech.net/gollum/api/v1/sep0024/transactions/withdraw/interactive", strings.NewReader(form.Encode()))
	if err != nil {
		return InteractiveResp{}, err
	}
	interactiveReq.Header.Add("Authorization", bearer)
	interactiveReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(interactiveReq)
	if err != nil {
		return InteractiveResp{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	res := InteractiveResp{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return InteractiveResp{}, err
	}
	return res, nil
}

func SignTransaction(Seed, SigningKey, XDR string) (string, error) {
	kp, err := keypair.Parse(Seed)
	if err != nil {
		return "", fmt.Errorf("can't parse sender key, error: %v", err)
	}

	domains := []string{"demo-wallet-server.stellar.org", "gollum-sep24.armenotech.net"}
	val, challengeFor, _, err := stellar.ReadChallengeTx(XDR, SigningKey, network.TestNetworkPassphrase, "gollum-sep24.armenotech.net", domains)
	if err != nil {
		return "", err
	}

	if challengeFor != kp.Address() {
		return "", nil
	}
	tx := val
	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full))
	if err != nil {
		return "", fmt.Errorf("can't sign transaction, error: %v", err)
	}
	txe, err := tx.Base64()
	if err != nil {
		return "", fmt.Errorf("can't convert to base 64, error: %v", err)
	}
	fmt.Println(txe)

	authXDR := map[string]interface{}{
		"transaction": txe,
	}
	authBody, err := json.Marshal(authXDR)
	if err != nil {
		return "", err
	}

	resultAuth, err := http.Post("https://gollum-sep24.armenotech.net/gollum/api/v1/sep0010", "application/json", bytes.NewReader(authBody))
	if err != nil {
		return "", err
	}

	token := struct {
		Token string `json:"token"`
	}{}
	if err := json.NewDecoder(resultAuth.Body).Decode(&token); err != nil {
		return "", err
	}

	fmt.Println(token.Token)

	return token.Token, nil
}

type Wallet struct {
	WalletAddress string `json:"wallet_address"`
	WalletSeed    string `json:"wallet_seed"`
	Blockchain    string `json:"blockchain"`
}

func (service ATWalletService) CreateWallet(externalId, merchantId string) (string, string, string, error) {
	pair, err := keypair.Random()
	if err != nil {
		return "", "", "", err
	}
	wallet := pair.Address()
	seed := pair.Seed()

	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + wallet)
	if err != nil {
		return "", "", "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}
	fmt.Println(string(body))

	clientStellar := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: pair.Address()}
	sourceAccount, err := clientStellar.AccountDetail(ar)

	asset, err := stellar.ParseAssetString("ATUSD:GBT4VVTDPCNA45MNWX5G6LUTLIEENSTUHDVXO2AQHAZ24KUZUPLPGJZH")
	if err != nil {
		return "", "", "", err
	}
	op := stellar.ChangeTrust{
		SourceAccount: wallet,
		Line:          asset.MustToChangeTrustAsset(),
		Limit:         stellar.MaxTrustlineLimit,
	}

	tx, err := stellar.NewTransaction(
		stellar.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []stellar.Operation{&op},
			BaseFee:              stellar.MinBaseFee,
			Preconditions:        stellar.Preconditions{TimeBounds: stellar.NewInfiniteTimeout()}, // Use a real timeout in production!

		},
	)

	// Sign the transaction
	tx, err = tx.Sign(network.TestNetworkPassphrase, pair)
	if err != nil {
		log.Fatalln(err)
	}

	// Get the base 64 encoded transaction envelope
	txe, err := tx.Base64()
	if err != nil {
		log.Fatalln(err)
	}

	// Send the transaction to the network
	_, err = clientStellar.SubmitTransactionXDR(txe)
	if err != nil {
		log.Fatalln(err)
	}
	return seed, wallet, wallet, nil
}

func (service ATWalletService) GetExistingWallet(externalId, merchantId, blockchain string) (string, string, error) {

	wallet := Wallet{Blockchain: blockchain}

	_, key, walletByte, err := service.store.GetByUser(merchantId, externalId)

	if err != nil && errors.Is(err, storage.ErrNotFound) {

		wallet.WalletSeed, wallet.WalletAddress, key, err = service.CreateWallet(merchantId, externalId)
		if err != nil {
			return "", "", err
		}

		wallet.Blockchain = blockchain
		value, err := json.Marshal(wallet)
		if err != nil {
			return "", "", err
		}

		_, err = service.store.Put(merchantId, externalId, key, value)
		if err != nil {
			return "", "", err
		}

	} else if err != nil {
		return "", "", err
	} else {
		err = json.Unmarshal(walletByte, &wallet)
		if err != nil {
			return "", "", err
		}
	}
	return wallet.WalletAddress, wallet.WalletSeed, err
}
