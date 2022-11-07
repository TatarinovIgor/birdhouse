package telegram

import (
	"birdhouse/modules/service"
	"crypto/rsa"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (receiver *TelegramService) InitRouter() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := receiver.bot.GetUpdatesChan(u)
	bot := receiver.bot
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			receiver.Router(update)
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
	fmt.Println("exit")
}

func NewTelegramService(privateKey *rsa.PrivateKey, bot *tgbotapi.BotAPI, atWallet *service.ATWalletService) *TelegramService {
	return &TelegramService{
		privateKey: privateKey,
		bot:        bot,
		atWallet:   atWallet,
	}
}
