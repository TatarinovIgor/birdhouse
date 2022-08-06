package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"

	"net/url"
)

type Payload struct {
	ExternalId string `json:"external_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

type ActivationBody struct {
	Platform string `json:"platform"`
	Type     string `json:"type"`
	Name     string `json:"Name"`
}

func CreateWalletBH(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jwtToken := r.URL.Query().Get("auth_key")
	merchId := r.URL.Query().Get("merch_id")
	URL, _ := url.Parse("https://wallet.rock-west.net/api/v1/wallet/application/" + merchId + "/sign-up")
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

	if res.StatusCode != 200 {
		w.WriteHeader(res.StatusCode)
		w.Write([]byte(res.Status))
		return
	}
	accessToken := gjson.Get(string(resBytes), "access_token")
	data := `{"platform":"stellar","type":"ccba7c71-27aa-40c3-9fe8-03db6934bc20","name":"private account"}`

	ActivationURL, _ := url.Parse("https://wallet.rock-west.net/api/v1/wallet/application/" + merchId + "/user/platform/stellar/account")
	ActivationReq, err := http.NewRequest("POST", ActivationURL.String(), bytes.NewBuffer([]byte(data)))
	ActivationReq.Header.Set("Content-Type", "text/application/json; charset=utf-8")
	ActivationReq.Header.Set("Authorization", "Bearer "+accessToken.Str)
	ActivationReq.Header.Set("X-Auth-Token", jwtToken)
	ActivationReq.Header.Set("X-Session-ID", sessionId.String())

	ActivationRes, err := http.DefaultClient.Do(ActivationReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`error forming request`))
		return
	}

	ActivationResBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`error parsing wallet response`))
		return
	}

	if ActivationRes.StatusCode == 200 {
		guid := gjson.Get(string(ActivationResBytes), "guid")

		w.WriteHeader(ActivationRes.StatusCode)
		w.Write([]byte(guid.Str))
		return
	}

	w.WriteHeader(ActivationRes.StatusCode)
	w.Write([]byte(ActivationRes.Status))
	return
}

func CreateWalletAT(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jwtToken := r.URL.Query().Get("X-Auth-Token")
	_ = r.URL.Query().Get("X-Session-ID")

	claims := jwt.MapClaims{}

	decoded, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not parse jwt"))
		return
	}
	data := Payload{}
	fmt.Printf(decoded.Raw)
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
		if key == "external_id" {
			data.ExternalId = fmt.Sprintf("%v", val)
		} else if key == "first_name" {
			data.FirstName = fmt.Sprintf("%v", val)
		} else if key == "last_name" {
			data.LastName = fmt.Sprintf("%v", val)
		} else if key == "email" {
			data.Email = fmt.Sprintf("%v", val)
		} else if key == "phone" {
			data.Phone = fmt.Sprintf("%v", val)
		}
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not return data"))
		return
	}

	w.WriteHeader(200)
	return
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
