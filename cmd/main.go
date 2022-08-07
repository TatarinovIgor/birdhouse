package main

import (
	"birdhouse/modules/routing"
	"birdhouse/modules/service"
	"fmt"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT env variable must be set")
	}
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		log.Fatal("$BASE_PATH env variable must be set")
	}
	appGUIDStr := os.Getenv("APP_GUID")
	if appGUIDStr == "" {
		log.Fatal("$APP_GUID env variable must be set")
	}
	appGUID, err := uuid.Parse(appGUIDStr)
	if err != nil {
		log.Fatalf("incorrect $APP_GUID env variable: %s, error: %v", appGUIDStr, err)
	}
	atWalletService := service.NewATWalletService(basePath, "", appGUID)
	router := httprouter.New()
	urlPath := ""
	fmt.Println("hello i am started")
	routing.InitRouter(router, urlPath, atWalletService)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		fmt.Println("error", err)
		return
	}

}
