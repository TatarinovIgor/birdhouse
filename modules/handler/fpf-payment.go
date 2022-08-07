package handler

import (
	"birdhouse/modules/service"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

func MakeFPFLinkForWallet(atWallet *service.ATWalletService, isDeposit bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.URL.Query().Get("auth_key")
		amount, _ := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
		account := r.URL.Query().Get("acc_guid")
		accountGUID, err := uuid.Parse(account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := atWallet.SignIn(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		payment, err := atWallet.FPFPayment(jwtToken, token.AccessToken,
			"ATUSD", "GBT4VVTDPCNA45MNWX5G6LUTLIEENSTUHDVXO2AQHAZ24KUZUPLPGJZH",
			accountGUID, amount, isDeposit)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
		err = json.NewEncoder(w).Encode(payment.Action.Action)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
