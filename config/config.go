package config

import (
	"github.com/spf13/viper"
	"log"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using default settings: %v", err)
	}

	// Установка значений по умолчанию
	viper.SetDefault("database.file", "tasks.json")
	viper.SetDefault("server.port", 8080)
}
