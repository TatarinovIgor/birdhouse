package middleware

import (
	"birdhouse/modules/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func AuthMiddleware(atWallet *service.ATWalletService, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authToken := r.URL.Query().Get("auth_key")
		_, err := atWallet.TokenDecode(authToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next(w, r.WithContext(r.Context()), ps)
	}
}
