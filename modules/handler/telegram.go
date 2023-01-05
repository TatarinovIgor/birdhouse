package handler

import (
	"birdhouse/modules/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func TgAuth(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		//return auth page
	}
}
