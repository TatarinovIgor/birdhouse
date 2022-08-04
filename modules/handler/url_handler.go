package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"

	"net/url"
)

func CreateWalletBH(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jwtToken := r.URL.Query().Get("auth_key")
	URL, _ := url.Parse("https://wallet.rock-west.net/api/v1/wallet/application/ab54ee14-15f1-4ce5-bcc3-6559451354da/sign-up")
	sessionId := uuid.New()

	req := &http.Request{
		Method: "POST",
		URL:    URL,
		Header: map[string][]string{
			"X-Auth-Token": {jwtToken},
			"X-Session-ID": {sessionId.String()},
		},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`error forming request`))
		return
	}

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`error parsing wallet response`))
		return
	}

	if res.StatusCode == 200 {
		_, err := json.Marshal(resBytes)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`error formatting json response`))
			return
		}

		w.WriteHeader(res.StatusCode)
		accessToken := gjson.Get(string(resBytes), "access_token")
		w.Write([]byte(accessToken.Str))
		return
	}

	w.WriteHeader(res.StatusCode)
	w.Write([]byte(res.Status))
	return
}

func CreateWalletAT(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func SignInWalletBH(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func SignInAT(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func GetBalance(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func Deposit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func Withdraw(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func Transfer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func ListOfTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
