package handler

import (
	"birdhouse/modules/service"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func GetTransactionList(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.URL.Query().Get("auth_key")
		account := r.URL.Query().Get("acc_guid")
		accountGUID, err := uuid.Parse(account)
		token, err := atWallet.SignIn(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, "could not parse request data", http.StatusBadRequest)
			return
		}
		res, err := atWallet.GetTransactions(jwtToken, token.AccessToken, accountGUID)
		if err != nil {
			log.Println(err)
			http.Error(w, "error getting transaction list", http.StatusFailedDependency)
			return
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println(err)
			http.Error(w, "error parsing the response", http.StatusInternalServerError)
			return
		}
	}
}
