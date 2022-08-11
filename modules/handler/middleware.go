package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		auth_token := r.URL.Query().Get("auth_key")

		jwtKey := "secret"

		claims := jwt.MapClaims{}
		tkn, err := jwt.ParseWithClaims(auth_token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			fmt.Println(err.Error())
			fmt.Println(tkn)

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if claims, ok := tkn.Claims.(jwt.MapClaims); ok {
			// obtains claims
			fmt.Println(strconv.FormatInt(time.Now().Unix()+300, 10))
			fmt.Printf("sub = %v", uint(claims["iat"].(float64)))
			if uint(claims["iat"].(float64)) < uint(time.Now().UTC().Unix()-300) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r.WithContext(r.Context()), ps)
	}
}
