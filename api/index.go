package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// هيكل البيانات القادمة من التطبيق المصغر
type WebAppSignal struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	Action   string `json:"action"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// 1. إعدادات CORS (ضرورية جداً ليعمل الـ Fetch من المتصفح)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// التعامل مع طلبات الـ Preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return
	}

	// 2. تحليل البيانات القادمة
	var rawData json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	// محاولة معالجة الطلب كـ "إشارة ترحيب" من التطبيق المصغر
	var signal WebAppSignal
	if err := json.Unmarshal(rawData, &signal); err == nil && signal.Action == "welcome_trigger" {
		// إرسال رسالة ترحيب فورية
		welcomeText := fmt.Sprintf("أهلاً بك يا %s! ✨ لقد دخلت الآن إلى الاختبار، بالتوفيق!", signal.UserName)
		bot.Send(tgbotapi.NewMessage(signal.UserID, welcomeText))

		// جلب وإرسال صورة البروفايل
		photos, err := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: signal.UserID, Limit: 1})
		if err == nil && photos.TotalCount > 0 {
			photoMsg := tgbotapi.NewPhoto(signal.UserID, tgbotapi.FileID(photos.Photos[0][0].FileID))
			photoMsg.Caption = "صورة بروفايلك منورة التطبيق! 📸"
			bot.Send(photoMsg)
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 3. معالجة الطلب كـ "Update" عادي من تيليجرام (Start / Callback)
	var update tgbotapi.Update
	if err := json.Unmarshal(rawData, &update); err == nil {
		handleTelegramUpdate(bot, update)
	}

	w.WriteHeader(http.StatusOK)
}

func handleTelegramUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	const channelUsername = "@boxtoolls"
	const directLink = "http://t.me/TinyTuneBot/visuals"

	// معالجة أزرار التحقق
	if update.CallbackQuery != nil {
		handleCallback(bot, update.CallbackQuery, channelUsername, directLink)
		return
	}

	// معالجة أمر /start
	if update.Message != nil && update.Message.Text == "/start" {
		handleStart(bot, update.Message, channelUsername, directLink)
	}
}

func handleCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, channel, link string) {
	if query.Data == "verify_sub" {
		member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{SuperGroupUsername: channel, UserID: query.From.ID},
		})

		if err == nil && (member.Status == "member" || member.Status == "administrator" || member.Status == "creator") {
			bot.Send(tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID))
			msg := tgbotapi.NewMessage(query.Message.Chat.ID, fmt.Sprintf("أهلاً بك يا %s! 🌟 اضغط بالأسفل للدخول.", query.From.FirstName))
			msg.ReplyMarkup = createInlineKeyboard("🔗 دخول الاختبار", link)
			bot.Send(msg)
		} else {
			bot.Request(tgbotapi.NewCallbackWithAlert(query.ID, "❌ يجب عليك الاشتراك في القناة أولاً!"))
		}
	}
}

func handleStart(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, channel, link string) {
	member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{SuperGroupUsername: channel, UserID: msg.From.ID},
	})

	if err == nil && (member.Status == "member" || member.Status == "administrator" || member.Status == "creator") {
		newMsg := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("مرحباً بك مجدداً يا %s 👋", msg.From.FirstName))
		newMsg.ReplyMarkup = createInlineKeyboard("✨ دخول الاختبار", link)
		bot.Send(newMsg)
	} else {
		newMsg := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("أهلاً بك يا %s! ⚠️ اشترك أولاً لتتمكن من الدخول.", msg.From.FirstName))
		btnSub := tgbotapi.NewInlineKeyboardButtonURL("📢 اشترك في القناة", "https://t.me/boxtoolls")
		btnVerify := tgbotapi.NewInlineKeyboardButtonData("✅ تحقق من الاشتراك", "verify_sub")
		newMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btnSub), tgbotapi.NewInlineKeyboardRow(btnVerify))
		bot.Send(newMsg)
	}
}

func createInlineKeyboard(text, url string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(text, url)))
}
