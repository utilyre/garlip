package config

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var defaultConfig Config

var (
	ErrNotSet = errors.New("not set")
)

type Config struct {
	DatabaseURL string
	TokenSecret []byte
}

func init() {
	var err error

	defaultConfig.DatabaseURL, err = lookupEnv("DB_URL")
	if err != nil {
		log.Fatal("config:", err)
	}

	tokenSecret, err := lookupEnv("JWT_SECRET")
	if err != nil {
		log.Fatal("config:", err)
	}
	defaultConfig.TokenSecret = []byte(tokenSecret)
}

func Default() *Config {
	return &defaultConfig
}

func lookupEnv(key string) (string, error) {
	value, exist := os.LookupEnv(key)
	if !exist {
		return "", fmt.Errorf("environment variable '%s': %w", key, ErrNotSet)
	}

	return value, nil
}
