package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	BotToken   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppEnv     string
	LogLevel   string
	Port       string
	RateLimit  int
	RateWindow string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	config := &Config{
		BotToken:   os.Getenv("BOT_TOKEN"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		AppEnv:     os.Getenv("APP_ENV"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
		Port:       os.Getenv("PORT"),
		RateLimit:  100,
		RateWindow: "1h",
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.BotToken == "" {
		return fmt.Errorf("BOT_TOKEN is required")
	}
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DBPort == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	return nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
