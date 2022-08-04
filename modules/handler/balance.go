package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func GetBalance(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	/*
		jwtToken := r.Header.Get("auth_key")
		guid := r.Header.Get("guid")

		//sign_in_data := sign_in_wallet_send_data(jwtToken)

		//header := {"Authorization": "Bearer " + sign_in_data["access_token"]}
		URL := "https://atwallet.rock-west.net/api/v1/wallet/application/ab54ee14-15f1-4ce5-bcc3-6559451354da/user/platform/stellar/account/"

		request, _ := http.NewRequest("POST", URL, nil)
		request.Header.Get(URL + guid + "?include_assets=true" + header)
		//return json.Unmarshal(request)

	*/
}
