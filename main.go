package main

import (
	"birdhouse/modules/routing"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	urlPath := ""

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := httprouter.New()

	fmt.Println("hello i am started")
	routing.InitRouter(router, urlPath)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		fmt.Println("error", err)
		return
	}

}
