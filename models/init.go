// Package models содержит логику для работы с клиентами и их лимитами в базе данных.
package models

import (
	"cloud/logger"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

// DB глобальная переменная для работы с базой данных PostgreSQL через GORM.
var DB *gorm.DB

// RDB глобальная переменная для работы с базой данных Redis.
var RDB *redis.Client

// InitENV загружает переменные окружения из файла .env.
// Завершает выполнение программы в случае ошибки загрузки.
func InitENV() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.PrintFatal("Ошибка при получении .env файла")
		return
	}
}

// InitDB инициализирует подключение к базе данных PostgreSQL.
// Проверяет наличие таблицы "rate_limits". В случае ошибок завершает выполнение программы.
func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.PrintFatal(err.Error())
	}

	logger.PrintInfo("Успешное подключение к БД")

	// Получаем список всех таблиц схемы public
	var tables []string
	result := DB.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)
	if result.Error != nil {
		logger.PrintFatal(result.Error.Error())
	}

	// Проверка наличия таблицы "rate_limits"
	found := false
	for _, table := range tables {
		if table == "rate_limits" {
			found = true
			break
		}
	}
	if !found {
		logger.PrintFatal("Отсутствует таблица \"rate_limits\"")
	}
}

// InitRDB инициализирует подключение к Redis-серверу.
// Выполняет проверку соединения с помощью команды PING.
// Завершает выполнение программы в случае ошибки.
func InitRDB() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := RDB.Ping(ctx).Result()
	if err != nil {
		logger.PrintFatal("Ошибка подключения к Redis: " + err.Error())
	}
	logger.PrintInfo("Успешное подключение к Redis: " + pong)
}
