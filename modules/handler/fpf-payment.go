package handler

import (
	"birdhouse/modules/service"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func MakeFPFLinkForWallet(atWallet *service.ATWalletService, isDeposit bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var DepositRequest service.DepositRequest
		err := json.NewDecoder(r.Body).Decode(&DepositRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if isDeposit {
			deposit, err := atWallet.CreateStellarDeposit(DepositRequest.ExternalID, DepositRequest.MerchantID, DepositRequest.Blockchain)
			if err != nil {
				log.Println(err)
				http.Error(w, "could not perform deposit", http.StatusBadRequest)
				return
			}
			err = json.NewEncoder(w).Encode(deposit)
			if err != nil {
				log.Println(err)
				http.Error(w, "could not parse response from server", http.StatusInternalServerError)
				return
			}
		} else {
			withdraw, err := atWallet.CreateStellarWithdraw(DepositRequest.ExternalID, DepositRequest.MerchantID, DepositRequest.Blockchain)
			if err != nil {
				log.Println(err)
				http.Error(w, "could not perform deposit", http.StatusBadRequest)
				return
			}
			err = json.NewEncoder(w).Encode(withdraw)
			if err != nil {
				log.Println(err)
				http.Error(w, "could not parse response from server", http.StatusInternalServerError)
				return
			}
		}
	}
}
