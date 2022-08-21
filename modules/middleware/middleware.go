package middleware

import (
	"birdhouse/modules/service"
	"context"
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

func WrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Take the context out from the request
		ctx := r.Context()

		// Get new context with key-value "params" -> "httprouter.Params"
		ctx = context.WithValue(ctx, "params", ps)

		// Get new http.Request with the new context
		r = r.WithContext(ctx)

		// Call your original http.Handler
		h.ServeHTTP(w, r)
	}
}
