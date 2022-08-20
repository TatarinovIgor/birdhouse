package handler

import (
	"birdhouse/modules/service"
	"fmt"
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
		_, err = fmt.Fprintf(w, "%s", deposit.StellarMemo)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
