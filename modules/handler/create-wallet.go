package handler

import (
	"birdhouse/modules/service"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func MakeCreateSignUPWithWalletBH(atWallet *service.ATWalletService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		jwtToken := r.URL.Query().Get("auth_key")
		token, err := atWallet.SignUp(jwtToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
		result, err := atWallet.CreateStellarWallet(jwtToken, token.AccessToken,
			"ccba7c71-27aa-40c3-9fe8-03db6934bc20", "BirdHouseClientAccount")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
		err = json.NewEncoder(w).Encode(result.GUID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
	}
}
