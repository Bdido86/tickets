package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const (
	apiToken = "TELEGRAM_BOT_API_TOKEN"
	debug    = "DEBUG"
)

var singleInstance *Config

type Config struct {
	token string
	debug bool
}

func (c Config) Token() string {
	return c.token
}

func (c Config) Debug() bool {
	return c.debug
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found.")
		panic("No .env file")
	}
}

func GetInstance() *Config {
	if singleInstance == nil {
		singleInstance = &Config{
			token: getEnv(apiToken, ""),
			debug: getEnvAsBool(debug, false),
		}
	}
	return singleInstance
}

func getEnv(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
