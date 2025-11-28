package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	WaterCount    int
	TrainingTime  int
	CardioTime    int
	StrengthTime  int
	TotalCalories int
	Protein       int
	Fat           int
	Carbs         int
	CurrentDay    int
	LastActivity  string
	LastCommand   string
}

var users = make(map[int64]*User)

func (u *User) cheсkDayUpdate() {
	today := time.Now().Format("2006-01-02")

	if u.LastActivity != today {
		u.WaterCount = 0
		u.TrainingTime = 0
		u.CardioTime = 0
		u.StrengthTime = 0
		u.TotalCalories = 0
		u.Protein = 0
		u.Fat = 0
		u.Carbs = 0
		u.CurrentDay++
		u.LastActivity = today

		log.Printf("🔄 Новый день пуся! День %d", u.CurrentDay)
	}
}

func getUser(chatID int64) *User {
	if users[chatID] == nil {
		users[chatID] = &User{}
	}
	return users[chatID]
}

func main() {
	loadUserData()

	bot, err := tgbotapi.NewBotAPI("8573098280:AAHtpTPlMpa2J3X5yLPOKJjcHgzepyvLnAY")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Твой Фит-Ботя готов к тренировкам !!! %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			user := getUser(callback.Message.Chat.ID)
			user.cheсkDayUpdate()

			log.Printf("[%s] нажал кнопку: %s", callback.From.UserName, callback.Data)
			switch callback.Data {
			case "cardio_15":
				user.CardioTime += 15
				user.TrainingTime += 15
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🏃 Добавлено 15 минут кардио! Отлично пуся! ❤️")
				bot.Send(msg)
				saveUserData()
			case "strength_15":
				user.StrengthTime += 15
				user.TrainingTime += 15
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "💪 Добавлено 15 минут силовой! Ты мощь пуся! 🔥")
				bot.Send(msg)
				saveUserData()
			}
			bot.Request(tgbotapi.NewCallback(callback.ID, ""))
			continue
		}
		if update.Message == nil {
			continue
		}

		user := getUser(update.Message.Chat.ID)
		user.cheсkDayUpdate()

		log.Printf("[%s], %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `🏋️ Привет! Я твой Фит-Ботя - лучший личный фитнес-помощник!

*Что я умею:*
/water - добавить воду 💧
/food - учет питания 🍎
/training - выбор тренировки ⏱️
/stats - статистика 📊

*Питание:*
/addcalories - добавить калории 🔥
/addprotein - добавить белки 🥩
/addfat - добавить жиры 🥑
/addcarbs - добавить углеводы 🍚

Давай начнем тренироваться вместе! 💪`)
			msg.ParseMode = "Markdown"
			bot.Send(msg)

		case "/water":
			user.WaterCount++
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "💧 Отлично! Вот столько выпито водички: "+strconv.Itoa(user.WaterCount))
			bot.Send(msg)
			saveUserData()

		case "/training":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🎯 *Выбери тип тренировки:*")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🏃 Кардио 15 мин", "cardio_15"),
					tgbotapi.NewInlineKeyboardButtonData("💪 Силовая 15 мин", "strength_15"),
				),
			)
			bot.Send(msg)

		case "/training15":
			user.TrainingTime += 15
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏱️ Отлично! Добавлено 15 минут тренировки, так держать пус! Всего: "+strconv.Itoa(user.TrainingTime)+" минут.")
			bot.Send(msg)
			saveUserData()

		case "/training30":
			user.TrainingTime += 30
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏱️ Отлично! Добавлено 30 минут тренировки, очень хорошо пус! Всего: "+strconv.Itoa(user.TrainingTime)+" минут.")
			bot.Send(msg)
			saveUserData()

		case "/training45":
			user.TrainingTime += 45
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏱️ Отлично! Добавлено 45 минут тренировки, молодчина пус! Всего: "+strconv.Itoa(user.TrainingTime)+" минут.")
			bot.Send(msg)
			saveUserData()

		case "/training60":
			user.TrainingTime += 60
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏱️ ВАААУ! целый час тренировки!!! Горжусь тобой пус! Всего: "+strconv.Itoa(user.TrainingTime)+" минут.")
			bot.Send(msg)
			saveUserData()

		case "/food", "🍎 Еда":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🍎 *Учет питания*\n\n"+
				"Выбери что хочешь добавить:\n"+
				"/addcalories - только калории\n"+
				"/addprotein - только белки\n"+
				"/addfat - только жиры\n"+
				"/addcarbs - только углеводы\n"+
				"/addall - все БЖУ сразу\n\n"+
				"Или введи данные вручную:\n"+
				"*250 20 10 30* - калории, белки, жиры, углеводы")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("🔥 Калории"),
					tgbotapi.NewKeyboardButton("🥩 Белки"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("🥑 Жиры"),
					tgbotapi.NewKeyboardButton("🍚 Углеводы"),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("📊 Все БЖУ"),
				),
			)
			bot.Send(msg)

		case "🔥 Калории", "/addcalories":
			user.LastCommand = "addcalories"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"🔥 Введи количество калорий:\nПример: *250*")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)

		case "🥩 Белки", "/addprotein":
			user.LastCommand = "addprotein"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"🥩 Введи количество белка (в граммах):\nПример: *25*")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)

		case "🥑 Жиры", "/addfat":
			user.LastCommand = "addfat"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"🥑 Введи количество жиров (в граммах):\nПример: *15*")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)

		case "🍚 Углеводы", "/addcarbs":
			user.LastCommand = "addcarbs"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"🍚 Введи количество углеводов (в граммах):\nПример: *40*")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)

		case "📊 Все БЖУ", "/addall":
			user.LastCommand = "addall"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"📊 Введи все данные через пробел:\n*Калории Белки Жиры Углеводы*\n\nПример: *250 20 10 30*")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)

		case "/stats":
			hours := user.TrainingTime / 60
			minutes := user.TrainingTime % 60
			statsText := "📊 Твоя статистика:\n" +
				"💧 Водичка: " + strconv.Itoa(user.WaterCount) + " стаканов\n" +
				"⏱️ Тренировки: " + strconv.Itoa(user.TrainingTime) + " минут"

			if hours > 0 {
				statsText += "\n🏆 Это " + strconv.Itoa(hours) + " часов " + strconv.Itoa(minutes) + " минут!"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, statsText)
			bot.Send(msg)

		default:
			// 1. Сначала проверяем ввод отдельных чисел для БЖУ
			if number, err := strconv.Atoi(update.Message.Text); err == nil {
				user := getUser(update.Message.Chat.ID)

				switch user.LastCommand {
				case "addcalories":
					user.TotalCalories += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"🔥 Добавлено *"+strconv.Itoa(number)+"* ккал\n"+
							"Всего за день: *"+strconv.Itoa(user.TotalCalories)+"* ккал")
					msg.ParseMode = "Markdown"
					bot.Send(msg)
					saveUserData()

				case "addprotein":
					user.Protein += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"🥩 Добавлено *"+strconv.Itoa(number)+"*г белка\n"+
							"Всего за день: *"+strconv.Itoa(user.Protein)+"*г")
					msg.ParseMode = "Markdown"
					bot.Send(msg)
					saveUserData()

				case "addfat":
					user.Fat += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"🥑 Добавлено *"+strconv.Itoa(number)+"*г жиров\n"+
							"Всего за день: *"+strconv.Itoa(user.Fat)+"*г")
					msg.ParseMode = "Markdown"
					bot.Send(msg)
					saveUserData()

				case "addcarbs":
					user.Carbs += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"🍚 Добавлено *"+strconv.Itoa(number)+"*г углеводов\n"+
							"Всего за день: *"+strconv.Itoa(user.Carbs)+"*г")
					msg.ParseMode = "Markdown"
					bot.Send(msg)
					saveUserData()

				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"Сначала выбери что добавить через /food")
					bot.Send(msg)
				}
				continue
			}

			parts := strings.Fields(update.Message.Text)
			if len(parts) == 4 {
				calories, err1 := strconv.Atoi(parts[0])
				protein, err2 := strconv.Atoi(parts[1])
				fat, err3 := strconv.Atoi(parts[2])
				carbs, err4 := strconv.Atoi(parts[3])

				if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
					user.TotalCalories += calories
					user.Protein += protein
					user.Fat += fat
					user.Carbs += carbs

					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"🍎 Добавлено:\n"+
							"🔥 "+strconv.Itoa(calories)+" ккал\n"+
							"🥩 "+strconv.Itoa(protein)+"г белка\n"+
							"🥑 "+strconv.Itoa(fat)+"г жиров\n"+
							"🍚 "+strconv.Itoa(carbs)+"г углеводов\n\n"+
							"Всего за день: "+strconv.Itoa(user.TotalCalories)+" ккал")
					bot.Send(msg)
					saveUserData()
					continue
				}
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"Используй команды:\n"+
					"/start - меню\n"+
					"/water - добавить воду\n"+
					"/training - выбрать тренировку\n"+
					"/food - добавить питание\n"+
					"/stats - статистика")
			bot.Send(msg)
		}
	}
}
