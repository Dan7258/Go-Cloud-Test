package configHandler

import (
	"cloud/logger"
	"encoding/json"
	"os"
)

type Config struct {
	Port          string   `json:"port"`
	Backends      []string `json:"backends"`
	Capacity      int      `json:"capacity"`
	RatePerSecond int      `json:"rate_per_second"`
}

func (c *Config) Init() {
	file, err := os.Open("config.json")
	if err != nil {
		logger.PrintFatal(err.Error())
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		logger.PrintFatal(err.Error())
	}
	if c.Backends == nil {
		logger.PrintFatal("список серверов пуст")
	}
	if c.Port == "" {
		logger.PrintFatal("отсутствует порт прослушивания")
	}
	logger.PrintInfo("Получены данные из config.json")
}
