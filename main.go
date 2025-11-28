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

// Главное меню
func getMainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("💧 Вода"),
			tgbotapi.NewKeyboardButton("🏋️ Тренировка"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🍎 Питание"),
			tgbotapi.NewKeyboardButton("📊 Статистика"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🧹 Очистить"),
		),
	)
}

// Меню тренировок
func getTrainingKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏃 Кардио 15", "cardio_15"),
			tgbotapi.NewInlineKeyboardButtonData("💪 Силовая 15", "strength_15"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏃 Кардио 30", "cardio_30"),
			tgbotapi.NewInlineKeyboardButtonData("💪 Силовая 30", "strength_30"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏃 Кардио 45", "cardio_45"),
			tgbotapi.NewInlineKeyboardButtonData("💪 Силовая 45", "strength_45"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏃 Кардио 60", "cardio_60"),
			tgbotapi.NewInlineKeyboardButtonData("💪 Силовая 60", "strength_60"),
		),
	)
}

// Меню питания
func getFoodKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔥 Калории", "calories"),
			tgbotapi.NewInlineKeyboardButtonData("🥩 Белки", "protein"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥑 Жиры", "fat"),
			tgbotapi.NewInlineKeyboardButtonData("🍚 Углеводы", "carbs"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Все БЖУ", "all_nutrients"),
		),
	)
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
		// Обработка нажатий на инлайн-кнопки
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
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "💪 Добавлено 15 минут силовой! Так держать пус! 🔥")
				bot.Send(msg)
				saveUserData()
			case "cardio_30":
				user.CardioTime += 30
				user.TrainingTime += 30
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🏃 Добавлено 30 минут кардио! Супер пуся! 🌟")
				bot.Send(msg)
				saveUserData()
			case "strength_30":
				user.StrengthTime += 30
				user.TrainingTime += 30
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "💪 Добавлено 30 минут силовой! Невероятно пус! 💥")
				bot.Send(msg)
				saveUserData()
			case "cardio_45":
				user.CardioTime += 45
				user.TrainingTime += 45
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🏃 Добавлено 45 минут кардио! Фантастика пуся! 🚀")
				bot.Send(msg)
				saveUserData()
			case "strength_45":
				user.StrengthTime += 45
				user.TrainingTime += 45
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "💪 Добавлено 45 минут силовой! Ты монстр пус! 🤯")
				bot.Send(msg)
				saveUserData()
			case "cardio_60":
				user.CardioTime += 60
				user.TrainingTime += 60
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🏃 ЦЕЛЫЙ ЧАС КАРДИО!!! Ты ИДЕАЛ пуся! 👑")
				bot.Send(msg)
				saveUserData()
			case "strength_60":
				user.StrengthTime += 60
				user.TrainingTime += 60
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "💪 ЦЕЛЫЙ ЧАС СИЛОВОЙ!!! Ты легенда пус! 🏆")
				bot.Send(msg)
				saveUserData()
			case "calories":
				user.LastCommand = "addcalories"
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🔥 Введи количество калорий заюсь:\nПример: *250*")
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			case "protein":
				user.LastCommand = "addprotein"
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🥩 Введи количество белка зай (в граммах):\nПример: *25*")
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			case "fat":
				user.LastCommand = "addfat"
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🥑 Введи количество жиров пуся (в граммах):\nПример: *15*")
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			case "carbs":
				user.LastCommand = "addcarbs"
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🍚 Введи количество углеводов пус (в граммах):\nПример: *40*")
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			case "all_nutrients":
				user.LastCommand = "addall"
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "📊 Введи все данные через пробел пуся:\n*Калории Белки Жиры Углеводы*\n\nПример: *250 20 10 30*")
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			}

			bot.Request(tgbotapi.NewCallback(callback.ID, ""))
			continue
		}

		// Обработка обычных сообщений
		if update.Message == nil {
			continue
		}

		user := getUser(update.Message.Chat.ID)
		user.cheсkDayUpdate()

		log.Printf("[%s], %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Text {
		case "/start", "/menu":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `🏋️ Привет пуся! Я твой Фит-Ботя - лучший личный фитнес-помощник!

Выбирай что будем делать сегодня! 💪`)
			msg.ReplyMarkup = getMainKeyboard()
			bot.Send(msg)

		case "💧 Вода":
			user.WaterCount++
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "💧 Отлично пуся! Выпито стаканов водички: "+strconv.Itoa(user.WaterCount))
			bot.Send(msg)
			saveUserData()

		case "🏋️ Тренировка":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🎯 *Выбери тип и продолжительность тренировки пус:*")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = getTrainingKeyboard()
			bot.Send(msg)

		case "🍎 Питание":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🍎 *Что добавим пус?*")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = getFoodKeyboard()
			bot.Send(msg)

		case "📊 Статистика":
			hours := user.TrainingTime / 60
			minutes := user.TrainingTime % 60
			statsText := "📊 *Твоя статистика пупся:*\n\n" +
				"💧 Водичка: " + strconv.Itoa(user.WaterCount) + " стаканов\n" +
				"⏱️ Тренировки: " + strconv.Itoa(user.TrainingTime) + " минут\n" +
				"🏃 Кардио: " + strconv.Itoa(user.CardioTime) + " минут\n" +
				"💪 Силовая: " + strconv.Itoa(user.StrengthTime) + " минут\n" +
				"🔥 Калории: " + strconv.Itoa(user.TotalCalories) + " ккал\n" +
				"🥩 Белки: " + strconv.Itoa(user.Protein) + "г\n" +
				"🥑 Жиры: " + strconv.Itoa(user.Fat) + "г\n" +
				"🍚 Углеводы: " + strconv.Itoa(user.Carbs) + "г"

			if hours > 0 {
				statsText += "\n\n🏆 *Это " + strconv.Itoa(hours) + " часов " + strconv.Itoa(minutes) + " минут тренировок!*"
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, statsText)
			msg.ParseMode = "Markdown"
			bot.Send(msg)

		case "🧹 Очистить":
			user.WaterCount = 0
			user.TrainingTime = 0
			user.CardioTime = 0
			user.StrengthTime = 0
			user.TotalCalories = 0
			user.Protein = 0
			user.Fat = 0
			user.Carbs = 0
			user.LastCommand = ""
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🧹 *Все данные очищены пус!*\n\n💧 Стаканы воды: 0\n⏱️ Время тренировок: 0 мин\n🔥 Калории: 0\n🥩 Белки: 0г\n🥑 Жиры: 0г\n🍚 Углеводы: 0г\n\nНачинаем с чистого листа заюсь! 💫")
			msg.ParseMode = "Markdown"
			bot.Send(msg)
			saveUserData()

		default:
			// Обработка ввода чисел для БЖУ
			if number, err := strconv.Atoi(update.Message.Text); err == nil {
				switch user.LastCommand {
				case "addcalories":
					user.TotalCalories += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🔥 Добавлено *"+strconv.Itoa(number)+"* ккал зай!\nВсего за день: *"+strconv.Itoa(user.TotalCalories)+"* ккал")
					msg.ParseMode = "Markdown"
					bot.Send(msg)
					saveUserData()

				case "addprotein":
					user.Protein += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🥩 Добавлено *"+strconv.Itoa(number)+"*г белка пус!\nВсего за день: *"+strconv.Itoa(user.Protein)+"*г")
					msg.ParseMode = "Markdown"
					bot.Send(msg)
					saveUserData()

				case "addfat":
					user.Fat += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🥑 Добавлено *"+strconv.Itoa(number)+"*г жиров зай!\nВсего за день: *"+strconv.Itoa(user.Fat)+"*г")
					msg.ParseMode = "Markdown"
					bot.Send(msg)
					saveUserData()

				case "addcarbs":
					user.Carbs += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🍚 Добавлено *"+strconv.Itoa(number)+"*г углеводов зай!\nВсего за день: *"+strconv.Itoa(user.Carbs)+"*г")
					msg.ParseMode = "Markdown"
					bot.Send(msg)
					saveUserData()

				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала выбери что добавить через меню '🍎 Питание' пуся")
					bot.Send(msg)
				}
				continue
			}

			// Обработка ввода всех БЖУ сразу
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
						"🍎 Добавлено пуся:\n"+
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

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Используй кнопки меню пус! 🎯")
			msg.ReplyMarkup = getMainKeyboard()
			bot.Send(msg)
		}
	}
}
