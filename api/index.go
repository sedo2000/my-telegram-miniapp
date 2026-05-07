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
		w.WriteHeader(http.StatusOK)
		return
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "TinyTune Bot: يعمل بنجاح 🚀")
		return
	}

	const channelUsername = "@boxtoolls"
	const directLink = "http://t.me/TinyTuneBot/visuals"

	// 1. معالجة الضغط على زر التحقق (Callback Query)
	if update.CallbackQuery != nil {
		chatID := update.CallbackQuery.Message.Chat.ID
		userID := update.CallbackQuery.From.ID
		firstName := update.CallbackQuery.From.FirstName // الحصول على الاسم هنا أيضاً

		if update.CallbackQuery.Data == "verify_sub" {
			member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
				ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
					SuperGroupUsername: channelUsername,
					UserID:             userID,
				},
			})

			if err == nil && (member.Status == "member" || member.Status == "administrator" || member.Status == "creator") {
				// حذف رسالة الاشتراك
				bot.Send(tgbotapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID))
				
				// رسالة الترحيب بعد التحقق
				welcomeMsg := fmt.Sprintf("أهلاً بك يا %s في اختبار التمويل! 🌟\n\nاضغط على الزر بالأسفل للدخول للأختبار.", firstName)
				msg := tgbotapi.NewMessage(chatID, welcomeMsg)
				
				button := map[string]interface{}{
					"text": "🔗 دخول الاختبار",
					"url":  directLink,
				}
				keyboard := map[string]interface{}{"inline_keyboard": [][]interface{}{{button}}}
				keyboardBytes, _ := json.Marshal(keyboard)
				msg.ReplyMarkup = json.RawMessage(keyboardBytes)
				bot.Send(msg)
			} else {
				// تنبيه في حال عدم الاشتراك
				callbackConfig := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "❌ عذراً، يجب عليك الاشتراك في القناة أولاً!")
				bot.Request(callbackConfig)
			}
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 2. معالجة أمر البداية /start
	if update.Message != nil && update.Message.Text == "/start" {
		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID
		firstName := update.Message.From.FirstName // الحصول على اسم المستخدم

		member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				SuperGroupUsername: channelUsername,
				UserID:             userID,
			},
		})

		if err == nil && (member.Status == "member" || member.Status == "administrator" || member.Status == "creator") {
			// المستخدم مشترك بالفعل
			welcomeMsg := fmt.Sprintf("مرحباً بك مجدداً يا %s في اختبار التمويل 👋", firstName)
			msg := tgbotapi.NewMessage(chatID, welcomeMsg)
			
			button := map[string]interface{}{"text": "✨ دخول الاختبار", "url": directLink}
			keyboard := map[string]interface{}{"inline_keyboard": [][]interface{}{{button}}}
			keyboardBytes, _ := json.Marshal(keyboard)
			msg.ReplyMarkup = json.RawMessage(keyboardBytes)
			bot.Send(msg)
		} else {
			// المستخدم غير مشترك - رسالة الاشتراك الإجباري
			welcomeMsg := fmt.Sprintf("أهلاً بك يا %s! ⚠️\n\nيجب عليك الاشتراك في قناة البوت أولاً لتتمكن من الدخول إلى اختبار التمويل.", firstName)
			msg := tgbotapi.NewMessage(chatID, welcomeMsg)
			
			btnSub := map[string]interface{}{"text": "📢 اشترك في القناة", "url": "https://t.me/boxtoolls"}
			btnVerify := map[string]interface{}{"text": "✅ تحقق من الاشتراك", "callback_data": "verify_sub"}
			
			keyboard := map[string]interface{}{
				"inline_keyboard": [][]interface{}{{btnSub}, {btnVerify}},
			}
			
			keyboardBytes, _ := json.Marshal(keyboard)
			msg.ReplyMarkup = json.RawMessage(keyboardBytes)
			bot.Send(msg)
		}
	}

	w.WriteHeader(http.StatusOK)
}
