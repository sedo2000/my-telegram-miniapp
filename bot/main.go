package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// 1. ضع التوكن الخاص بك هنا بدلاً من YOUR_BOT_TOKEN
	bot, err := tgbotapi.NewBotAPI("8458116007:AAHU-Ch47PVdOJOH8LmzPL_UXxAwQrTHUlQ")
	if err != nil {
		log.Panicf("خطأ في التوكن: %v", err)
	}

	// تفعيل الـ Debug لرؤية كل شيء في الـ Terminal
	bot.Debug = true
	log.Printf("تم تشغيل البوت بنجاح: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// تجاهل أي تحديث ليس رسالة نصية
		if update.Message == nil {
			continue
		}

		// طباعة الرسالة الواردة في الـ Terminal للتأكد من وصولها
		log.Printf("وصلت رسالة من [%s]: %s", update.Message.From.UserName, update.Message.Text)

		// التحقق من أمر /start أو النص الصريح /start
		if update.Message.IsCommand() && update.Message.Command() == "start" || update.Message.Text == "/start" {
			
			// رابط الـ Mini App (تأكد من صحته من BotFather)
			targetURL := "https://t.me/TinyTuneBot/visuals"

			// إنشاء زر الرابط المضمون
			button := tgbotapi.NewInlineKeyboardButtonURL("🎵 Open TinyTune Visuals", targetURL)
			
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(button),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to **TinyTune**! 🚀\n\nالبوت شغال الآن، اضغط على الزر للبدء:")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = keyboard

			// إرسال الرسالة مع فحص الخطأ
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("فشل في إرسال الرد: %v", err)
			} else {
				log.Println("تم إرسال الرد بنجاح!")
			}
		}
	}
}
