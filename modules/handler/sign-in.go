package handler

import (
	"birdhouse/modules/service"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func MakeSignInWalletBH(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.Header.Get("auth_key")
		token, err := atWallet.SignUp(jwtToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
		err = json.NewEncoder(w).Encode(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
	}
}

func MakeSignInAT(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.Header.Get("X-Auth-Token")
		//sessionID := r.Header.Get("X-Session-ID")
		result, err := atWallet.TokenDecode(jwtToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
	}
}
