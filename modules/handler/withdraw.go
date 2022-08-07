package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Withdraw(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	/*
		jwtToken := r.Header.Get("auth_key")
		guid := r.Header.Get("guid")
		amount := r.Header.Get("amount")
	*/
}
