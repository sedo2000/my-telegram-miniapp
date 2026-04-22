package main

import (
	"log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// 1. ضع توكن البوت الخاص بك هنا
	bot, err := tgbotapi.NewBotAPI("8458116007:AAHU-Ch47PVdOJOH8LmzPL_UXxAwQrTHUlQ")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("TinyTune Bot Started: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		if update.Message.Command() == "start" {
			// 2. رابط الـ Mini App (رابط Vercel الخاص بك)
			webApp := &tgbotapi.WebAppInfo{
				URL: "https://my-telegram-miniapp.vercel.app/",
			}

			// إنشاء زر الـ Mini App
			button := tgbotapi.NewInlineKeyboardButtonWebApp("🎵 Launch TinyTune Visuals", *webApp)
			
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(button),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to **TinyTune**! 🚀\n\nExperience high-end visuals and music.")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = keyboard

			bot.Send(msg)
		}
	}
}
