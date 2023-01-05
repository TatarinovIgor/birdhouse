package telegram_service

import (
	"birdhouse/modules/service"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

type TelegramService struct {
	bot         *tgbotapi.BotAPI
	atWallet    *service.ATWalletService
	currentUser *service.UsersData
	walletUrl   string
	walletKey   string
}

var menu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Main-menu"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Pay in"),
		tgbotapi.NewKeyboardButton("Pay out"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Balance"),
	),
)

var ButtonMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Sign in", "https://no-code-wallet.bird-house.org/home"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Transaction list", "1"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Transfer", "2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Rates", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("About", "google.com"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Terms and conditions", "google.com"),
	),
)

func (receiver *TelegramService) ListenAndServe() {
	bot := receiver.bot
	var err error
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, I am BirdHouse bot")
				msg.ReplyMarkup = menu
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
			case "Main-menu":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Menu:")
				msg.ReplyMarkup = ButtonMenu
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
			case "Link-Wallet":
				receiver.LinkTelegram(update)
			case "Deposit":
				receiver.Deposit(update)
			case "Withdraw":
				receiver.Withdraw(update)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Menu")
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
	fmt.Println("exit")
}
func NewTelegramService(bot *tgbotapi.BotAPI, atWallet *service.ATWalletService, walletUrl string, walletKey string) *TelegramService {
	return &TelegramService{
		bot:       bot,
		atWallet:  atWallet,
		walletUrl: walletUrl,
		walletKey: walletKey,
	}
}

func (receiver TelegramService) LinkTelegram(update tgbotapi.Update) {
	url := receiver.walletUrl + "/link-telegram/?phone=" + strings.Split(update.Message.Text, " ")[len(strings.Split(update.Message.Text, " "))] + "&telegram_token=" + strconv.FormatInt(update.Message.From.ID, 10) + "&api_token=" + receiver.walletKey
	_, err := receiver.atWallet.TelegramRequester(url)
	if err != nil {
		panic(err)
	}
	print(update.Message.Text)
}

func (receiver TelegramService) Deposit(update tgbotapi.Update) {
	_, err := receiver.atWallet.GetUsersFromDatabase()
	if err != nil {
		panic(err)
	}
	print(update.Message.Text)
}

func (receiver TelegramService) Withdraw(update tgbotapi.Update) {
	_, err := receiver.atWallet.GetUsersFromDatabase()
	if err != nil {
		panic(err)
	}
	print(update.Message.Text)
}
