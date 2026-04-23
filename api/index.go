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
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return
	}

	if update.Message != nil && update.Message.Text == "/start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to TinyTune! 🚀")
		
		// بناء هيكل الزر يدوياً كـ Map لضمان نجاح الـ Build
		// r.Host سيجلب رابط مشروعك على فيرسل تلقائياً
		webAppURL := "https://" + r.Host + "/"

		// إنشاء الـ Keyboard باستخدام Raw JSON
		// هذه الطريقة تتخطى خطأ "undefined: WebAppInfo" نهائياً
		button := map[string]interface{}{
			"text": "🎵 Open TinyTune Visualizer",
			"web_app": map[string]string{
				"url": webAppURL,
			},
		}

		keyboard := map[string]interface{}{
			"inline_keyboard": [][]interface{}{
				{button},
			},
		}

		keyboardBytes, _ := json.Marshal(keyboard)
		msg.ReplyMarkup = json.RawMessage(keyboardBytes)

		bot.Send(msg)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
