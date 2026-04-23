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
	log.Printf("TinyTune is LIVE on: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		if update.Message.Command() == "start" {
			// استخدم رابط الـ Mini App الرسمي الذي حصلت عليه من BotFather
			// أو رابط Vercel مباشرة
			targetURL := "http://t.me/TinyTuneBot/visuals"

			// إنشاء زر رابط عادي - هذا السطر مستحيل يعطيك خطأ
			button := tgbotapi.NewInlineKeyboardButtonURL("🎵 Open TinyTune Visuals", targetURL)
			
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(button),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to **TinyTune**! 🚀\n\nClick the button to start.")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = keyboard

			bot.Send(msg)
		}
	}
}
