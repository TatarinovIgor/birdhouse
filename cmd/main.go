package main

import (
	"birdhouse/modules/routing"
	"birdhouse/modules/service"
	telegram_service "birdhouse/modules/telegram-service"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/oklog/run"
	"log"
	"net/http"
	"os"
	"strconv"
)

// @title Example API
// @version 1.0
// @description Description of the example API
// @securityDefinitions.basic BasicAut

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
	publicKeyPath := os.Getenv("PUBLIC_KEY")
	if appGUIDStr == "" {
		log.Fatal("$PUBLIC_KEY env variable must be set")
	}
	tokenTimeToLiveStr := os.Getenv("TOKEN_TIME_TO_LIVE")
	if tokenTimeToLiveStr == "" {
		log.Fatal("$TOKEN_TIME_TO_LIVE env variable must be set")
	}
	tokenTimeToLive, err := strconv.ParseInt(tokenTimeToLiveStr, 10, 64)
	if err != nil {
		log.Fatal("could not convert to int $TOKEN_TIME_TO_LIVE")
	}
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatalf("could not read public key: %s, error: %v", publicKey, err)
	}
	block, _ := pem.Decode(publicKey)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("could not parse public key env variable: %s, error: %v", publicKey, err)
	}
	appGUID, err := uuid.Parse(appGUIDStr)
	if err != nil {
		log.Fatalf("incorrect $APP_GUID env variable: %s, error: %v", appGUIDStr, err)
	}
	bot, err := tgbotapi.NewBotAPI("5698836967:AAEO1kCse9XP5xDw67RYWOs9tSsZHpDlFDM")
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	atWalletService := service.NewATWalletService(basePath, pub, appGUID, tokenTimeToLive)
	telegramService := telegram_service.NewTelegramService(bot, updates, atWalletService)

	router := httprouter.New()
	urlPath := ""
	fmt.Println("hello i am started")

	routing.InitRouter(router, urlPath, atWalletService)

	g := run.Group{}
	// stream manager
	{
		g.Add(func() error {
			fmt.Println("telegram service starting")
			telegramService.ListenAndServe()
			return nil
		}, func(err error) {
			fmt.Println("telegram service stopping")
		})
	}
	// REST API
	{
		g.Add(func() error {
			fmt.Println("REST API starting")
			err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
			if err != nil {
				fmt.Println("error", err)
				return err
			}
			return nil
		}, func(err error) {
			fmt.Println("REST API stopping")
		})
	}
	fmt.Println("app exiting")

	// documentation for share
	// opts1 := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	// sh1 := middleware.Redoc(opts1, nil)
	// r.Handle("/docs", sh1)
}
