package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// читаем переменные окружения из .env
func LoadEnv() {
	err := godotenv.Load("app/.env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
}

// собираем строку для подключения к БД
func GetDBConnString() string {
	return "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"
}

// передаем secret-key
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET не задан явно в .env")
	}
	return secret
}
