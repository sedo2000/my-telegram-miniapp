package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// 1. ضع التوكن الخاص بك هنا
	bot, err := tgbotapi.NewBotAPI("8458116007:AAHU-Ch47PVdOJOH8LmzPL_UXxAwQrTHUlQ")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("TinyTune Bot is LIVE: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		if update.Message.Command() == "start" {
			// 2. رابط الـ Mini App الخاص بك (الذي حصلت عليه من BotFather)
			// الرابط يكون بهذا الشكل: https://t.me/TinyTuneBot/visuals
			miniAppURL := "http://t.me/TinyTuneBot/visuals"

			// إنشاء زر رابط عادي (هذا سيعمل حتماً)
			button := tgbotapi.NewInlineKeyboardButtonURL("🎵 Launch TinyTune Visuals", miniAppURL)
			
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(button),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to **TinyTune**! 🚀\n\nClick below to start the experience.")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = keyboard

			bot.Send(msg)
		}
	}
}
