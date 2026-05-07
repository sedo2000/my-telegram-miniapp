package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// هيكل استقبال البيانات من التطبيق المصغر (الـ Fetch الذي وضعته في HTML)
type WebAppSignal struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	Action   string `json:"action"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	// قراءة البيانات القادمة في الطلب
	var bodyBytes []byte
	if r.Body != nil {
		var err error
		// سنحاول فك التشفير مرتين، لذا نحتاج لقراءة البيانات أولاً
		decoder := json.NewDecoder(r.Body)
		
		// محاولة 1: هل الطلب "إشارة ترحيب" من التطبيق المصغر؟
		var signal WebAppSignal
		// نقوم بنسخ محتوى الـ Body لنتمكن من استخدامه مرة أخرى إذا فشل التحويل الأول
		// لتبسيط الأمر، سنحاول فك تشفير البيانات مباشرة كـ Signal
		err = json.Unmarshal(bodyBytes, &signal) // ملاحظة: هذا للتبسيط المنطقي
	}

	// --- الجزء الأول: معالجة إشارة الترحيب من التطبيق المصغر ---
	// سنقوم بتحويل الـ Body وفحصه
	var rawData map[string]interface{}
	json.NewDecoder(r.Body).Decode(&rawData)

	if action, ok := rawData["action"].(string); ok && action == "welcome_trigger" {
		userID := int64(rawData["user_id"].(float64))
		userName := rawData["user_name"].(string)

		// 1. إرسال ترحيب نصي
		welcomeText := fmt.Sprintf("أهلاً بك يا %s! ✨ لقد دخلت الآن إلى الاختبار، بالتوفيق!", userName)
		bot.Send(tgbotapi.NewMessage(userID, welcomeText))

		// 2. جلب صورة البروفايل الحقيقية وإرسالها
		photos, err := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: userID, Limit: 1})
		if err == nil && photos.TotalCount > 0 {
			fileID := photos.Photos[0][0].FileID
			photoMsg := tgbotapi.NewPhoto(userID, tgbotapi.FileID(fileID))
			photoMsg.Caption = "صورة بروفايلك منورة التطبيق! 📸"
			bot.Send(photoMsg)
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// --- الجزء الثاني: معالجة رسائل تيليجرام العادية (Start / Callback) ---
	// نقوم بتحويل البيانات إلى Update الخاص بـ Telegram
	var update tgbotapi.Update
	// ملاحظة: في بيئة الإنتاج يفضل استخدام مفسر واحد، هنا قمنا بدمج المنطق
	data, _ := json.Marshal(rawData)
	json.Unmarshal(data, &update)

	const channelUsername = "@boxtoolls"
	const directLink = "http://t.me/TinyTuneBot/visuals"

	// معالجة الضغط على زر التحقق
	if update.CallbackQuery != nil {
		chatID := update.CallbackQuery.Message.Chat.ID
		userID := update.CallbackQuery.From.ID
		firstName := update.CallbackQuery.From.FirstName

		if update.CallbackQuery.Data == "verify_sub" {
			member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
				ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
					SuperGroupUsername: channelUsername,
					UserID:             userID,
				},
			})

			if err == nil && (member.Status == "member" || member.Status == "administrator" || member.Status == "creator") {
				bot.Send(tgbotapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID))
				welcomeMsg := fmt.Sprintf("أهلاً بك يا %s في اختبار التمويل! 🌟\n\nاضغط على الزر بالأسفل للدخول للأختبار.", firstName)
				msg := tgbotapi.NewMessage(chatID, welcomeMsg)
				
				button := map[string]interface{}{"text": "🔗 دخول الاختبار", "url": directLink}
				keyboard := map[string]interface{}{"inline_keyboard": [][]interface{}{{button}}}
				kbBytes, _ := json.Marshal(keyboard)
				msg.ReplyMarkup = json.RawMessage(kbBytes)
				bot.Send(msg)
			} else {
				bot.Request(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "❌ عذراً، يجب عليك الاشتراك في القناة أولاً!"))
			}
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// معالجة أمر /start
	if update.Message != nil && update.Message.Text == "/start" {
		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID
		firstName := update.Message.From.FirstName

		member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				SuperGroupUsername: channelUsername,
				UserID:             userID,
			},
		})

		if err == nil && (member.Status == "member" || member.Status == "administrator" || member.Status == "creator") {
			welcomeMsg := fmt.Sprintf("مرحباً بك مجدداً يا %s في اختبار التمويل 👋", firstName)
			msg := tgbotapi.NewMessage(chatID, welcomeMsg)
			button := map[string]interface{}{"text": "✨ دخول الاختبار", "url": directLink}
			keyboard := map[string]interface{}{"inline_keyboard": [][]interface{}{{button}}}
			kbBytes, _ := json.Marshal(keyboard)
			msg.ReplyMarkup = json.RawMessage(kbBytes)
			bot.Send(msg)
		} else {
			welcomeMsg := fmt.Sprintf("أهلاً بك يا %s! ⚠️\n\nيجب عليك الاشتراك في قناة البوت أولاً لتتمكن من الدخول إلى اختبار التمويل.", firstName)
			msg := tgbotapi.NewMessage(chatID, welcomeMsg)
			btnSub := map[string]interface{}{"text": "📢 اشترك في القناة", "url": "https://t.me/boxtoolls"}
			btnVerify := map[string]interface{}{"text": "✅ تحقق من الاشتراك", "callback_data": "verify_sub"}
			keyboard := map[string]interface{}{"inline_keyboard": [][]interface{}{{btnSub}, {btnVerify}}}
			kbBytes, _ := json.Marshal(keyboard)
			msg.ReplyMarkup = json.RawMessage(kbBytes)
			bot.Send(msg)
		}
	}

	w.WriteHeader(http.StatusOK)
}
