package handler

import (
	"encoding/json"
	"net/http"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	bot, _ := tgbotapi.NewBotAPI("8458116007:AAHU-Ch47PVdOJOH8LmzPL_UXxAwQrTHUlQ")

	var update tgbotapi.Update
	json.NewDecoder(r.Body).Decode(&update)

	if update.Message != nil && update.Message.Text == "/start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to TinyTune via Webhook! 🚀")
		bot.Send(msg)
	}
}
