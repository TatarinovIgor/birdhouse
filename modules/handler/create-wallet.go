package handler

import (
	"birdhouse/modules/service"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// swagger:route GET /create_wallet_bh create new wallet
// Create new wallet
//
// security:
// - apiKey: []
// responses:
//  403: Forbidden
//  424: Failed Dependency
//  200: Success
func MakeCreateSignUPWithWalletBH(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Println("new signup req")
		jwtToken := r.URL.Query().Get("auth_key")
		token, err := atWallet.SignUp(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		result, err := atWallet.CreateStellarWallet(jwtToken, token.AccessToken,
			"ccba7c71-27aa-40c3-9fe8-03db6934bc20", "BirdHouseClientAccount")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
		_, err = fmt.Fprintf(w, "%s", result.GUID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
	}
}
