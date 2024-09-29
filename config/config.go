package config

import (
	"github.com/spf13/viper"
	"log"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Config file not found, using default settings")
	}

	// Установка значений по умолчанию
	viper.SetDefault("database.file", "tasks.json")
	viper.SetDefault("server.port", 8080)
}
