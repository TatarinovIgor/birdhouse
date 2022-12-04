package handler

import (
	"birdhouse/modules/service"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func MakeSignInWalletBH(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.URL.Query().Get("auth_key")
		token, err := atWallet.SignIn(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, "could not parse request data", http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(token)
		if err != nil {
			log.Println(err)
			http.Error(w, "could not parse response from server", http.StatusInternalServerError)
			return
		}
	}
}

func MakeSignInAT(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.Header.Get("X-Auth-Token")
		//sessionID := r.Header.Get("X-Session-ID")
		result, err := atWallet.TokenDecode(jwtToken)
		fmt.Println(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, "error decoding token", http.StatusFailedDependency)
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Println(err)
			http.Error(w, "error encoding the response", http.StatusInternalServerError)
			return
		}
	}
}
