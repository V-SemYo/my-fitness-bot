# my-fitness-bot 🏋️

Telegram-бот для трекинга питания и тренировок. Помогает вести учёт спортивных достижений и рациона прямо в мессенджере.

---

## Описание / Description

**RU**: Простой Telegram-бот на Go, который позволяет пользователям записывать свои тренировки и приёмы пищи, отслеживать прогресс и получать статистику.

**EN**: A simple Telegram fitness bot written in Go. Allows users to log workouts and meals, track progress, and view statistics.

---

## Функционал / Features

- 📝 Добавление записей о тренировках (вид, длительность, сожжённые калории)
- 🍎 Ведение дневника питания (продукты, порции, калории)
- 📊 Просмотр статистики за день / неделю
- 💾 Локальное хранение данных (in-memory storage)

---

## Технологии / Tech Stack

- **Go** 1.21+
- **Telegram Bot API** — взаимодействие с пользователем
- **In-memory storage** — хранение данных в оперативной памяти

---

## Установка и запуск / Installation

```bash
# Клонирование репозитория
git clone https://github.com/V-SemYo/my-fitness-bot.git
cd my-fitness-bot

# Установка зависимостей
go mod tidy

# Запуск (требуется переменная окружения с токеном бота)
export TELEGRAM_BOT_TOKEN="ваш_токен"
go run main.go storage.go
