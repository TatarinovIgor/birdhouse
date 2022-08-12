package middleware

import (
	"birdhouse/modules/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func AuthMiddleware(atWallet *service.ATWalletService, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		/*
			authToken := r.URL.Query().Get("auth_key")
				_, issuedAt, expiresAt, err := atWallet.TokenDecode(authToken)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if expiresAt.Before(time.Now()) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				fmt.Print(issuedAt, expiresAt)
		*/
		next(w, r.WithContext(r.Context()), ps)
	}
}
