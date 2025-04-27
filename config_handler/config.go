package configHandler

import (
	"cloud/logger"
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port     string   `json:"port"`
	Backends []string `json:"backends"`
}

func (c *Config) Init() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	if c.Backends == nil {
		logger.PrintFatal("список серверов пуст")
	}
	if c.Port == "" {
		logger.PrintFatal("отсутствует порт прослушивания")
	}
	logger.PrintInfo("Получены данные для балансировщика нагрузки")
}
