package main

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	WaterCount    float64
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
	WaterStep     float64
	Unit          string
	Streak        int
	LastStreakDay string
}

var users = make(map[int64]*User)
var botMessages = make(map[int64][]int)

// Фразы поддержки со смайликами
var supportPhrases = []string{
	"Отлично пуся! 💖",
	"Так держать заюсь! 🌟",
	"Молодец мася! 🥰",
	"Прекрасно пус! 💫",
	"Умничка зая! 🌈",
	"Великолепно мась! 🎀",
	"Супер пупся! 💕",
	"Замечательно зайка! 🌸",
}

// Генератор случайных фраз поддержки
func getSupportPhrase() string {
	return supportPhrases[time.Now().Unix()%int64(len(supportPhrases))]
}

func (u *User) checkDayUpdate() {
	today := time.Now().Format("2006-01-02")

	// Обновление стрика дней
	if u.LastStreakDay == "" {
		u.LastStreakDay = today
		u.Streak = 1
	} else {
		lastDay, _ := time.Parse("2006-01-02", u.LastStreakDay)
		currentDay, _ := time.Parse("2006-01-02", today)
		daysDiff := int(currentDay.Sub(lastDay).Hours() / 24)

		if daysDiff == 1 {
			// Последовательные дни
			u.Streak++
			u.LastStreakDay = today
		} else if daysDiff > 1 {
			// Пропущен день - сбрасываем стрик
			u.Streak = 1
			u.LastStreakDay = today
		}
	}

	// Обновление ежедневных данных
	if u.LastActivity != today {
		oldDay := u.CurrentDay
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
		if oldDay > 0 {
			log.Printf("🔄 Новый день пуся! День %d", u.CurrentDay)
		}
	}
}

func getUser(chatID int64) *User {
	if users[chatID] == nil {
		users[chatID] = &User{
			WaterStep: 0.2,
			Unit:      "л",
			Streak:    1,
		}
	}
	return users[chatID]
}

func deleteBotMessages(bot *tgbotapi.BotAPI, chatID int64) {
	log.Printf("🔍 ОТЛАДКА: Удаляю сообщения для chatID %d", chatID)
	if messages, exists := botMessages[chatID]; exists {
		log.Printf("🔍 ОТЛАДКА: Найдено %d сообщений для удаления", len(messages))
		for i, msgID := range messages {
			log.Printf("🔍 ОТЛАДКА: Удаляю сообщение %d: ID %d", i+1, msgID)
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, msgID)
			_, err := bot.Send(deleteMsg)
			if err != nil {
				log.Printf("❌ Ошибка удаления сообщения %d: %v", msgID, err)
			} else {
				log.Printf("✅ Успешно удалено сообщение %d", msgID)
			}
			time.Sleep(100 * time.Millisecond)
		}
		botMessages[chatID] = []int{}
	} else {
		log.Printf("❌ ОТЛАДКА: Нет сообщений для chatID %d", chatID)
	}
}

func addBotMessage(chatID int64, MessageID int) {
	log.Printf("➕ Добавляю сообщение %d для chatID %d", MessageID, chatID)
	botMessages[chatID] = append(botMessages[chatID], MessageID)
	if len(botMessages[chatID]) > 50 {
		botMessages[chatID] = botMessages[chatID][:50]
	}
	log.Printf("📊 Теперь в chatID %d: %d сообщений", chatID, len(botMessages[chatID]))
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
			tgbotapi.NewKeyboardButton("👤 Мой профиль"),
			tgbotapi.NewKeyboardButton("🧹 Очистить"),
		),
	)
}

// Inline-кнопки для воды
func getWaterKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("+0.2", "water_0.2"),
			tgbotapi.NewInlineKeyboardButtonData("+0.5", "water_0.5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("-0.2", "water_-0.2"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", "water_cancel"),
		),
	)
}

// Inline-кнопки для настроек воды
func getWaterSettingsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0.1", "step_0.1"),
			tgbotapi.NewInlineKeyboardButtonData("0.2", "step_0.2"),
			tgbotapi.NewInlineKeyboardButtonData("0.25", "step_0.25"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0.33", "step_0.33"),
			tgbotapi.NewInlineKeyboardButtonData("0.5", "step_0.5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("мл", "unit_ml"),
			tgbotapi.NewInlineKeyboardButtonData("л", "unit_l"),
			tgbotapi.NewInlineKeyboardButtonData("стаканы", "unit_glass"),
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

// Функция для создания progress bar
func getProgressBar(current, total float64, width int) string {
	if total == 0 {
		return strings.Repeat("░", width)
	}
	percentage := current / total
	filled := int(math.Round(percentage * float64(width)))
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

// Функция для отображения профиля пользователя
func getUserProfile(user *User) string {
	// Прогресс воды (предположим, цель 2 литра)
	waterGoal := 2.0
	if user.Unit == "мл" {
		waterGoal = 2000
	} else if user.Unit == "стаканы" {
		waterGoal = 8
	}

	waterProgress := getProgressBar(user.WaterCount, waterGoal, 10)
	waterPercentage := int((user.WaterCount / waterGoal) * 100)
	if waterPercentage > 100 {
		waterPercentage = 100
	}

	profile := "👤 *Твой профиль пуся!*\n\n"
	profile += "📅 *День:* " + strconv.Itoa(user.CurrentDay) + "\n"
	profile += "🔥 *Стрик дней:* " + strconv.Itoa(user.Streak) + " дней подряд!\n\n"

	profile += "💧 *Вода сегодня:* " + formatWater(user.WaterCount, user.Unit) + "\n"
	profile += "📊 " + waterProgress + " " + strconv.Itoa(waterPercentage) + "%\n\n"

	profile += "⚙️ *Настройки:*\n"
	profile += "Шаг воды: " + strconv.FormatFloat(user.WaterStep, 'f', -1, 64) + " " + user.Unit + "\n"
	profile += "Единицы: " + user.Unit + "\n\n"

	profile += getSupportPhrase()

	return profile
}

// Функция для форматирования воды в зависимости от единиц
func formatWater(amount float64, unit string) string {
	switch unit {
	case "мл":
		return strconv.FormatFloat(amount, 'f', 0, 64) + " мл"
	case "стаканы":
		return strconv.FormatFloat(amount, 'f', 1, 64) + " стаканов"
	default:
		return strconv.FormatFloat(amount, 'f', 1, 64) + " л"
	}
}

// Функция конвертации между единицами
func convertWater(value float64, fromUnit, toUnit string) float64 {
	if fromUnit == toUnit {
		return value
	}

	// Конвертируем в литры сначала
	var liters float64
	switch fromUnit {
	case "мл":
		liters = value / 1000
	case "стаканы":
		liters = value * 0.25 // 1 стакан = 250 мл = 0.25 л
	default:
		liters = value
	}

	// Конвертируем из литров в целевую единицу
	switch toUnit {
	case "мл":
		return liters * 1000
	case "стаканы":
		return liters * 4
	default:
		return liters
	}
}

func main() {
	loadUserData()

	bot, err := tgbotapi.NewBotAPI("8573098280:AAHtpTPlMpa2J3X5yLPOKJjcHgzepyvLnAY")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Твой Фит-Ботя готов к тренировкам! %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// Обработка нажатий на инлайн-кнопки
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			user := getUser(callback.Message.Chat.ID)
			user.checkDayUpdate()

			log.Printf("[%s] нажал кнопку: %s", callback.From.UserName, callback.Data)

			// Обработка кнопок воды
			if strings.HasPrefix(callback.Data, "water_") {
				parts := strings.Split(callback.Data, "_")
				if len(parts) == 2 {
					switch parts[1] {
					case "0.2", "0.5":
						value, _ := strconv.ParseFloat(parts[1], 64)
						convertedValue := convertWater(value, "л", user.Unit)
						user.WaterCount += convertedValue
						if user.WaterCount < 0 {
							user.WaterCount = 0
						}

						msg := tgbotapi.NewMessage(callback.Message.Chat.ID,
							"💧 "+getSupportPhrase()+"\n"+
								"Воды выпито: "+formatWater(user.WaterCount, user.Unit))
						if sentMsg, err := bot.Send(msg); err == nil {
							addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
						}
						saveUserData()

					case "-0.2":
						convertedValue := convertWater(0.2, "л", user.Unit)
						user.WaterCount -= convertedValue
						if user.WaterCount < 0 {
							user.WaterCount = 0
						}

						msg := tgbotapi.NewMessage(callback.Message.Chat.ID,
							"💧 Исправлено зай!\n"+
								"Воды выпито: "+formatWater(user.WaterCount, user.Unit))
						if sentMsg, err := bot.Send(msg); err == nil {
							addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
						}
						saveUserData()

					case "cancel":
						msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "❌ Отменено пус!")
						if sentMsg, err := bot.Send(msg); err == nil {
							addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
						}
					}
				}
			} else if strings.HasPrefix(callback.Data, "step_") {
				// Изменение шага воды
				step, _ := strconv.ParseFloat(strings.TrimPrefix(callback.Data, "step_"), 64)
				user.WaterStep = step

				msg := tgbotapi.NewMessage(callback.Message.Chat.ID,
					"⚙️ Шаг воды изменен на: "+strconv.FormatFloat(step, 'f', -1, 64)+" л\n"+getSupportPhrase())
				if sentMsg, err := bot.Send(msg); err == nil {
					addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
				}
				saveUserData()

			} else if strings.HasPrefix(callback.Data, "unit_") {
				// Изменение единиц измерения
				oldUnit := user.Unit
				user.Unit = strings.TrimPrefix(callback.Data, "unit_")

				// Конвертируем текущее значение воды в новые единицы
				user.WaterCount = convertWater(user.WaterCount, oldUnit, user.Unit)

				msg := tgbotapi.NewMessage(callback.Message.Chat.ID,
					"⚙️ Единицы изменены на: "+user.Unit+"\n"+getSupportPhrase())
				if sentMsg, err := bot.Send(msg); err == nil {
					addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
				}
				saveUserData()

			} else {
				// Обработка остальных callback'ов (тренировки, питание)
				switch callback.Data {
				case "cardio_15":
					user.CardioTime += 15
					user.TrainingTime += 15
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🏃 Добавлено 15 минут кардио! "+getSupportPhrase()+" ❤️")
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				case "strength_15":
					user.StrengthTime += 15
					user.TrainingTime += 15
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "💪 Добавлено 15 минут силовой! "+getSupportPhrase()+" 🔥")
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				case "cardio_30":
					user.CardioTime += 30
					user.TrainingTime += 30
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🏃 Добавлено 30 минут кардио! "+getSupportPhrase()+" 🌟")
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				case "strength_30":
					user.StrengthTime += 30
					user.TrainingTime += 30
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "💪 Добавлено 30 минут силовой! "+getSupportPhrase()+" 💥")
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				case "cardio_45":
					user.CardioTime += 45
					user.TrainingTime += 45
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🏃 Добавлено 45 минут кардио! "+getSupportPhrase()+" 🚀")
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				case "strength_45":
					user.StrengthTime += 45
					user.TrainingTime += 45
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "💪 Добавлено 45 минут силовой! "+getSupportPhrase()+" 🤯")
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				case "cardio_60":
					user.CardioTime += 60
					user.TrainingTime += 60
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🏃 ЦЕЛЫЙ ЧАС КАРДИО!!! "+getSupportPhrase()+" 👑")
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				case "strength_60":
					user.StrengthTime += 60
					user.TrainingTime += 60
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "💪 ЦЕЛЫЙ ЧАС СИЛОВОЙ!!! "+getSupportPhrase()+" 🏆")
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				case "calories":
					user.LastCommand = "addcalories"
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🔥 Введи количество калорий заюсь:\nПример: *250*")
					msg.ParseMode = "Markdown"
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
				case "protein":
					user.LastCommand = "addprotein"
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🥩 Введи количество белка зай (в граммах):\nПример: *25*")
					msg.ParseMode = "Markdown"
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
				case "fat":
					user.LastCommand = "addfat"
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🥑 Введи количество жиров пуся (в граммах):\nПример: *15*")
					msg.ParseMode = "Markdown"
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
				case "carbs":
					user.LastCommand = "addcarbs"
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🍚 Введи количество углеводов пус (в граммах):\nПример: *40*")
					msg.ParseMode = "Markdown"
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
				case "all_nutrients":
					user.LastCommand = "addall"
					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "📊 Введи все данные через пробел пуся:\n*Калории Белки Жиры Углеводы*\n\nПример: *250 20 10 30*")
					msg.ParseMode = "Markdown"
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(callback.Message.Chat.ID, sentMsg.MessageID)
					}
				}
			}

			bot.Request(tgbotapi.NewCallback(callback.ID, ""))
			continue
		}

		// Обработка обычных сообщений
		if update.Message == nil {
			continue
		}

		user := getUser(update.Message.Chat.ID)
		user.checkDayUpdate()

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Text {
		case "/start", "/menu":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🏋️ Привет пуся! Я твой Фит-Ботя - лучший личный фитнес-помощник! Выбирай что будем делать сегодня! 💪")
			msg.ReplyMarkup = getMainKeyboard()
			if sentMsg, err := bot.Send(msg); err == nil {
				addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
			}

		case "💧 Вода":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"💧 *Управление водой заюсь!*\n\n"+
					"Текущее количество: "+formatWater(user.WaterCount, user.Unit)+"\n"+
					"Шаг: "+strconv.FormatFloat(user.WaterStep, 'f', -1, 64)+" "+user.Unit+"\n\n"+
					"Выбери действие:")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = getWaterKeyboard()
			if sentMsg, err := bot.Send(msg); err == nil {
				addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
			}

		case "👤 Мой профиль":
			profile := getUserProfile(user)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, profile)
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = getWaterSettingsKeyboard()
			if sentMsg, err := bot.Send(msg); err == nil {
				addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
			}

		case "🏋️ Тренировка":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🎯 *Выбери тип и продолжительность тренировки пус:*")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = getTrainingKeyboard()
			if sentMsg, err := bot.Send(msg); err == nil {
				addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
			}

		case "🍎 Питание":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🍎 *Что добавим пус?*")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = getFoodKeyboard()
			if sentMsg, err := bot.Send(msg); err == nil {
				addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
			}

		case "📊 Статистика":
			hours := user.TrainingTime / 60
			minutes := user.TrainingTime % 60
			statsText := "📊 *Твоя статистика пупся:*\n\n" +
				"💧 Водичка: " + formatWater(user.WaterCount, user.Unit) + "\n" +
				"⏱️ Тренировки: " + strconv.Itoa(user.TrainingTime) + " минут\n" +
				"🏃 Кардио: " + strconv.Itoa(user.CardioTime) + " минут\n" +
				"💪 Силовая: " + strconv.Itoa(user.StrengthTime) + " минут\n" +
				"🔥 Калории: " + strconv.Itoa(user.TotalCalories) + " ккал\n" +
				"🥩 Белки: " + strconv.Itoa(user.Protein) + "г\n" +
				"🥑 Жиры: " + strconv.Itoa(user.Fat) + "г\n" +
				"🍚 Углеводы: " + strconv.Itoa(user.Carbs) + "г\n" +
				"🔥 Стрик: " + strconv.Itoa(user.Streak) + " дней подряд!"

			if hours > 0 {
				statsText += "\n\n🏆 *Это " + strconv.Itoa(hours) + " часов " + strconv.Itoa(minutes) + " минут тренировок!*"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, statsText)
			msg.ParseMode = "Markdown"
			if sentMsg, err := bot.Send(msg); err == nil {
				addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
			}

		case "🧹 Очистить":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🧹 Удаляю сообщения...")
			noticeMsg, _ := bot.Send(msg)
			if noticeMsg.MessageID != 0 {
				addBotMessage(update.Message.Chat.ID, noticeMsg.MessageID)
			}
			user.WaterCount = 0
			user.TrainingTime = 0
			user.CardioTime = 0
			user.StrengthTime = 0
			user.TotalCalories = 0
			user.Protein = 0
			user.Fat = 0
			user.Carbs = 0
			user.LastCommand = ""
			deleteBotMessages(bot, update.Message.Chat.ID)
			time.Sleep(1 * time.Second)
			if noticeMsg.MessageID != 0 {
				deleteNotice := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, noticeMsg.MessageID)
				bot.Send(deleteNotice)
			}
			clearMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "🧹 *Все данные очищены пус!*\n\n💧 Вода: 0\n⏱️ Время тренировок: 0 мин\n🔥 Калории: 0\n🥩 Белки: 0г\n🥑 Жиры: 0г\n🍚 Углеводы: 0г\n\nНачинаем с чистого листа заюсь! 💫")
			clearMsg.ParseMode = "Markdown"
			bot.Send(clearMsg)
			saveUserData()

		default:
			// Обработка ввода чисел для БЖУ
			if number, err := strconv.Atoi(update.Message.Text); err == nil {
				switch user.LastCommand {
				case "addcalories":
					user.TotalCalories += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🔥 Добавлено *"+strconv.Itoa(number)+"* ккал зай!\nВсего за день: *"+strconv.Itoa(user.TotalCalories)+"* ккал")
					msg.ParseMode = "Markdown"
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				case "addprotein":
					user.Protein += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🥩 Добавлено *"+strconv.Itoa(number)+"*г белка пус!\nВсего за день: *"+strconv.Itoa(user.Protein)+"*г")
					msg.ParseMode = "Markdown"
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				case "addfat":
					user.Fat += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🥑 Добавлено *"+strconv.Itoa(number)+"*г жиров зай!\nВсего за день: *"+strconv.Itoa(user.Fat)+"*г")
					msg.ParseMode = "Markdown"
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				case "addcarbs":
					user.Carbs += number
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🍚 Добавлено *"+strconv.Itoa(number)+"*г углеводов зай!\nВсего за день: *"+strconv.Itoa(user.Carbs)+"*г")
					msg.ParseMode = "Markdown"
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала выбери что добавить через меню '🍎 Питание' пуся")
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
					}
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
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🍎 Добавлено пуся:\n"+
						"🔥 "+strconv.Itoa(calories)+" ккал\n"+
						"🥩 "+strconv.Itoa(protein)+"г белка\n"+
						"🥑 "+strconv.Itoa(fat)+"г жиров\n"+
						"🍚 "+strconv.Itoa(carbs)+"г углеводов\n\n"+
						"Всего за день: "+strconv.Itoa(user.TotalCalories)+" ккал")
					if sentMsg, err := bot.Send(msg); err == nil {
						addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
					}
					saveUserData()
					continue
				}
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Используй кнопки меню пус! 🎯")
			msg.ReplyMarkup = getMainKeyboard()
			if sentMsg, err := bot.Send(msg); err == nil {
				addBotMessage(update.Message.Chat.ID, sentMsg.MessageID)
			}
		}
	}
}
