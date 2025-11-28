package main

import (
	"encoding/json"
	"log"
	"os"
)

func saveUserData() {
	data, err := json.Marshal(users)
	if err != nil {
		log.Printf("Ошибка сохранения: %v", err)
		return
	}
	err = os.WriteFile("user_data.json", data, 0644)
	if err != nil {
		log.Printf("Ошибка записи файла: %v", err)
		return
	}
	log.Println("✅ Данные сохранены!")
}

func loadUserData() {
	data, err := os.ReadFile("user_data.json")
	if err != nil {
		log.Println("📝 Файл данных не найден, начинаем с чистого листа")
		return
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Printf("Ошибка загрузки: %v", err)
		return
	}
	log.Printf("✅ Данные загружены! Пользователей: %d", len(users))
}
