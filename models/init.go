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

var DB *gorm.DB
var RDB *redis.Client

func InitENV() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.PrintFatal("Ошибка при получении .env файла")
		return
	}
}

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

	// Получаем список таблиц
	var tables []string
	result := DB.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)
	if result.Error != nil {
		logger.PrintWarning(result.Error.Error())
	}
}

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
		logger.PrintWarning("Error pinging redis: " + err.Error())
	}
	logger.PrintInfo("Успешное подключение к Reids: " + pong)
}
