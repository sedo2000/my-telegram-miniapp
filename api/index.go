package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return
	}

	var update tgbotapi.Update
	json.NewDecoder(r.Body).Decode(&update)

	if update.Message != nil && update.Message.Text == "/start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to TinyTune! 🚀")
		
		// الطريقة اليدوية لإنشاء لوحة المفاتيح لتجنب أخطاء الإصدارات
		// الرابط هو رابط الـ Frontend الخاص بك على فيرسل
		url := "https://" + r.Host + "/" 

		// إنشاء زر WebApp باستخدام ميزة "Map" لضمان عدم حدوث خطأ أثناء الـ Build
		button := map[string]interface{}{
			"text": "🎵 Open TinyTune",
			"web_app": map[string]string{
				"url": url,
			},
		}

		keyboard := map[string]interface{}{
			"inline_keyboard": [][]interface{}{
				{button},
			},
		}

		data, _ := json.Marshal(keyboard)
		msg.ReplyMarkup = json.RawMessage(data)

		bot.Send(msg)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
