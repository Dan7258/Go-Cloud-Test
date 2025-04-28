// Package models содержит логику для работы с клиентами и их лимитами в базе данных.
package models

import (
	"context"
	"fmt"
	"time"
)

// SetDataInRedis сохраняет данные в Redis с указанным ключом и временем жизни.
// Возвращает ошибку, если операция не удалась.
func SetDataInRedis(key string, value []byte, timeLive time.Duration) error {
	// Создание контекста с таймаутом для операции
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Сохранение данных в Redis с указанием времени жизни (TTL)
	err := RDB.Set(ctx, fmt.Sprintf("%s", key), value, timeLive).Err()
	return err
}

// GetDataFromRedis извлекает данные из Redis по ключу.
// Возвращает данные в виде массива байтов и ошибку, если операция не удалась.
func GetDataFromRedis(key string) ([]byte, error) {
	// Создание контекста с таймаутом для операции
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Извлечение данных из Redis
	return RDB.Get(ctx, key).Bytes()
}

// GetAllKeysFromRedis получает все ключи из Redis, соответствующие шаблону "*"
// (все ключи в базе данных).
// Возвращает список всех ключей и ошибку, если операция не удалась.
func GetAllKeysFromRedis() ([]string, error) {
	// Создание контекста с таймаутом для операции
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение всех ключей из Redis
	return RDB.Keys(ctx, "*").Result()
}
