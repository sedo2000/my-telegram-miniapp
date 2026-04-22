import telebot
from telebot import types

# ضع التوكن الخاص بك هنا من BotFather
API_TOKEN = '8458116007:AAHU-Ch47PVdOJOH8LmzPL_UXxAwQrTHUlQ'
# ضع رابط موقعك بعد رفعه على فيرسل هنا
MINI_APP_URL = 'https://your-project-name.vercel.app'

bot = telebot.TeleBot(API_TOKEN)

@bot.message_handler(commands=['start'])
def start(message):
    markup = types.InlineKeyboardMarkup()
    # إنشاء زر شفاف يفتح الـ Mini App
    web_app = types.WebAppInfo(MINI_APP_URL)
    button = types.InlineKeyboardButton(text="🎵 Open Music Experience", web_app=web_app)
    markup.add(button)
    
    bot.send_message(message.chat.id, "مرحباً بك! اضغط على الزر أدناه لتجربة الميني آب الجديد:", reply_markup=markup)

print("البوت يعمل الآن...")
bot.polling()
