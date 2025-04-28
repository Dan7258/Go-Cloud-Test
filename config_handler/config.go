// Package configHandler предоставляет функционал для загрузки и инициализации конфигурации приложения из JSON-файла.
package configHandler

import (
	"cloud/logger"
	"encoding/json"
	"os"
)

// Config представляет структуру конфигурационного файла.
// Содержит настройки порта, серверов-бэкендов, пропускной способности и скорости запросов.
type Config struct {
	Port       string   `json:"port"`         // Порт, на котором будет запущен сервер.
	Backends   []string `json:"backends"`     // Список адресов серверов-бэкендов.
	Capacity   int      `json:"capacity"`     // Емкость очереди запросов.
	RatePerSec int      `json:"rate_per_sec"` // Количество запросов в секунду.
}

// Init инициализирует конфигурацию, загружая её из файла config.json.
// При ошибке чтения или некорректных данных завершает выполнение программы.
func (c *Config) Init() {
	file, err := os.Open("config.json")
	if err != nil {
		logger.PrintFatal(err.Error()) // Логирует фатальную ошибку при открытии файла.
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		logger.PrintFatal(err.Error()) // Логирует фатальную ошибку при декодировании JSON.
	}

	// Проверка обязательных полей конфигурации
	if c.Backends == nil {
		logger.PrintFatal("список серверов пуст")
	}
	if c.Port == "" {
		logger.PrintFatal("отсутствует порт прослушивания")
	}

	logger.PrintInfo("Получены данные из config.json") // Лог успешной загрузки конфигурации.
}
