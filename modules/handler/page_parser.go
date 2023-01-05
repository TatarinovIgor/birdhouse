package handler

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

type PageVariables struct {
	Date string
	Time string
}

func CheckoutPage() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		HomePageVars := PageVariables{}
		t, err := template.ParseFiles("./templates/checkout.html")
		if err != nil { // if there is an error
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"` + `template parsing error` + `"}`))
			return
		}
		err = t.Execute(w, HomePageVars)
		if err != nil { // if there is an error
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"` + `template parsing error` + `"}`))
			return
		}
	}
}
