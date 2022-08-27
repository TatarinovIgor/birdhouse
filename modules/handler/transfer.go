package handler

import (
	"birdhouse/modules/service"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

func TransferDeposit(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.URL.Query().Get("auth_key")
		amount, _ := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
		accGuid := r.URL.Query().Get("acc_guid")
		accountReceiverExternalId := r.URL.Query().Get("receiver_external_id")
		accountSenderInternalId := r.URL.Query().Get("sender_internal_id")
		token, err := atWallet.SignIn(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		deposit, err := atWallet.Deposit(jwtToken, token.AccessToken,
			"ATUSD", "GBT4VVTDPCNA45MNWX5G6LUTLIEENSTUHDVXO2AQHAZ24KUZUPLPGJZH",
			accGuid, accountReceiverExternalId, accountSenderInternalId, amount)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}

		err = json.NewEncoder(w).Encode(deposit)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func TransferWithdraw(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.URL.Query().Get("auth_key")
		amount, _ := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
		accGuid := r.URL.Query().Get("acc_guid")
		accountSenderExternalId := r.URL.Query().Get("sender_external_id")
		accountReceiverInternalId := r.URL.Query().Get("receiver_internal_id")
		token, err := atWallet.SignIn(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		deposit, err := atWallet.Withdraw(jwtToken, token.AccessToken,
			"ATUSD", "GBT4VVTDPCNA45MNWX5G6LUTLIEENSTUHDVXO2AQHAZ24KUZUPLPGJZH",
			accGuid, accountSenderExternalId, accountReceiverInternalId, amount)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}

		err = json.NewEncoder(w).Encode(deposit)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
