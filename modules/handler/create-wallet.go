package handler

import (
	"birdhouse/modules/service"
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func MakeCreateWalletBH(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.Header.Get("auth_key")
		token, err := atWallet.SignUp(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
		result, err := atWallet.Activate(token.AccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
	}
}

func MakeCreateWalletAT(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.Header.Get("X-Auth-Token")
		//sessionID := r.Header.Get("X-Session-ID")
		result, err := atWallet.TokenDecode(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
	}
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
