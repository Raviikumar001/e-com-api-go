// internal/config/config.go

package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		Username string
		Password string
		Name     string
	}
	JWT struct {
		SecretKey string
		Issuer    string
	}
}

func LoadConfig() Config {
	return Config{
		Database: struct {
			Host     string
			Port     int
			Username string
			Password string
			Name     string
		}{
			Host:     os.Getenv("DB_HOST"),
			Port:     mustConvert(os.Getenv("DB_PORT")),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		JWT: struct {
			SecretKey string
			Issuer    string
		}{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
			Issuer:    os.Getenv("JWT_ISSUER"),
		},
	}
}

func mustConvert(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
