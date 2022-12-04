package handler

import (
	"birdhouse/modules/service"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// swagger:route GET /get_balance get user balance
// Get user balance
//
// security:
// - apiKey: []
// responses:
//  403: Forbidden
//  424: Failed Dependency
//  500: Internal Error
//  200: Success
func MakeGetBalance(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.URL.Query().Get("auth_key")
		token, err := atWallet.SignIn(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, "could not parse request data", http.StatusBadRequest)
			return
		}
		res, err := atWallet.GetBalance(jwtToken, token.AccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, "could not get current balance", http.StatusFailedDependency)
			return
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println(err)
			http.Error(w, "could not parse response from server", http.StatusInternalServerError)
			return
		}
	}
}
