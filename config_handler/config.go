package configHandler

import (
	"encoding/json"
	"errors"
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
		log.Fatal("список серверов пуст")
		panic(errors.New("список серверов пуст"))
	}
	if c.Port == "" {
		log.Fatal("отсутствует порт прослушивания")
		panic(errors.New("отсутствует порт прослушивания"))
	}
	log.Println("Получены данные для балансировщика нагрузки")

}
