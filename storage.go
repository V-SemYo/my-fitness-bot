package main

import (
	"encoding/json"
	"log"
	"os"
)

// Функции для сохранения и загрузки данных в файл
func saveUserData() {
	file, err := os.Create("userdata.json")
	if err != nil {
		log.Printf("❌ Ошибка создания файла: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(users); err != nil {
		log.Printf("❌ Ошибка кодирования: %v", err)
		return
	}
	log.Printf("💾 Данные сохранены для %d пользователей", len(users))
}

func loadUserData() {
	file, err := os.Open("userdata.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("📝 Файл данных не найден, начинаем с чистого листа")
			return
		}
		log.Printf("❌ Ошибка открытия файла: %v", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		log.Printf("❌ Ошибка декодирования данных: %v", err)
		return
	}
	log.Printf("📂 Данные загружены для %d пользователей", len(users))
}
