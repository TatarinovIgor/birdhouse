package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (receiver *TelegramService) Router(update tgbotapi.Update) {
	switch update.Message.Text {
	case "/start":
		receiver.startNewUser(update)
	case "Main menu":
		receiver.mainMenu(update)
	default:
		receiver.defaultCase(update)
	}
}
