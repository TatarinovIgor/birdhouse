package handler

import (
	"birdhouse/modules/service"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func KYC(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.URL.Query().Get("auth_key")
		accGuid := r.URL.Query().Get("acc_guid")
		token, err := atWallet.SignIn(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, "could not parse request data", http.StatusBadRequest)
			return
		}
		KYCServer, err := atWallet.KYCGetServer(jwtToken, token.AccessToken, accGuid)
		if err != nil {
			log.Println(err)
			http.Error(w, "error while getting KYC server", http.StatusFailedDependency)
			return
		}
		_ = fmt.Sprintf("%s/customer?id=%s", KYCServer, accGuid)

	}
}
