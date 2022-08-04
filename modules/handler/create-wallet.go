package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ATToken struct {
	AccessToken string `json:"access_token"'`
}

func CreateWalletBH(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jwtToken := r.Header.Get("auth_key")

	URL := "https://atwallet.rock-west.net/api/v1/wallet/application/ab54ee14-15f1-4ce5-bcc3-6559451354da/sign-up"
	// generate session and jwt from uid
	sessionID, _ := uuid.NewUUID()
	// create headers
	client := http.Client{}
	request, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	request.Header.Set("X-Auth-Token", jwtToken)
	request.Header.Set("X-Session-ID", sessionID.String())
	// make a request
	result, err := client.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}
	if result.StatusCode != 200 {
		http.Error(w, fmt.Sprintf("can't sign-up at wallet"), result.StatusCode)
		return
	}
	token := ATToken{}
	err = json.NewDecoder(result.Body).Decode(&token)
	if err != nil {
		http.Error(w, fmt.Sprintf("can't parse token from wallet"), result.StatusCode)
		return
	}
	activateWallet(token.AccessToken, jwtToken, sessionID.String())
}

func CreateWalletAT(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		jwtToken := r.Header.Get("X-Auth-Token")
		sessionID := r.Header.Get("X-Session-ID")

		decoded, err := .(jwtToken, jws.General)
		//return jsonify - ?


}
func activateWallet(accessToken, jwt, sessionID string) ([]byte, error) {
	URL := "https://atwallet.rock-west.net/api/v1/wallet/application/ab54ee14-15f1-4ce5-bcc3-6559451354da/sign-in"
	body, _ := json.Marshal(map[string]interface{}{
		"platform": "stellar",
		"type":     "ccba7c71-27aa-40c3-9fe8-03db6934bc20",
		"name":     "Private account",
	})
	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "text/application/json; charset=utf-8")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("X-Auth-Token", jwt)
	request.Header.Set("X-Auth-Token", sessionID)

	client := http.Client{}

	result, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(result.Body).Decode(&body)

	return body, err
}
