package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (receiver *TelegramService) Router(message tgbotapi.Update) {
	switch message.Message.Text {
	case "/start":
		receiver.startNewUser(message)
	case "Main menu":
		receiver.mainMenu(message)
	default:
		receiver.defaultCase(message)
	}

}
