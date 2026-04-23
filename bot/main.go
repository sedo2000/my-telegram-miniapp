package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// 1. ضع التوكن الخاص بك هنا بدلاً من YOUR_BOT_TOKEN
	bot, err := tgbotapi.NewBotAPI("8458116007:AAHU-Ch47PVdOJOH8LmzPL_UXxAwQrTHUlQ")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("TinyTune is LIVE: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		if update.Message.Command() == "start" {
			// استخدم الرابط الرسمي للـ Mini App الخاص بك
			// هذا الرابط سيفتح الـ Web App تلقائياً داخل تلجرام
			targetURL := "http://t.me/TinyTuneBot/visuals"

			// إنشاء زر رابط تقليدي - مضمون العمل 100%
			button := tgbotapi.NewInlineKeyboardButtonURL("🎵 Open TinyTune Visuals", targetURL)
			
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(button),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to **TinyTune**! 🚀\n\nClick the button below to start the experience.")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = keyboard

			bot.Send(msg)
		}
	}
}
