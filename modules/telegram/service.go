package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (receiver *TelegramService) startNewUser(update tgbotapi.Update) {
	MsgEmail := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter your email (format: test@gmail.com)")
	_, err := receiver.bot.Send(MsgEmail)
	messageEmail := update.Message.Text
	MsgPhone := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter your phone (format: +123456789)")
	_, err = receiver.bot.Send(MsgPhone)
	messagePhone := update.Message.Text
	JWTToken, err := receiver.CreateJWT(update.Message, messageEmail, messagePhone)
	if err != nil {
		_ = tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown error appeared")
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, JWTToken.GUID.String())
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, I am BirdHouse bot")
	msg.ReplyMarkup = menu
	if _, err = receiver.bot.Send(msg); err != nil {
		panic(err)
	}
}

func (receiver *TelegramService) mainMenu(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Menu:")
	msg.ReplyMarkup = ButtonMenu
	if _, err := receiver.bot.Send(msg); err != nil {
		panic(err)
	}
}

func (receiver *TelegramService) defaultCase(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Menu")
	if _, err := receiver.bot.Send(msg); err != nil {
		panic(err)
	}
}
