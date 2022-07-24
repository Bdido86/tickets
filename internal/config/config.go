package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const (
	tgToken = "TELEGRAM_BOT_API_TOKEN"
	port    = "PORT"
	debug   = "DEBUG"
)

var singleInstance *Config

type Config struct {
	token string
	port  string
	debug bool
}

func (c Config) Token() string {
	return c.token
}

func (c Config) Port() string {
	return c.port
}

func (c Config) Debug() bool {
	return c.debug
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found.")
		panic("No .env file")
	}

	singleInstance = new()
}

func GetConfig() *Config {
	return singleInstance
}

func new() *Config {
	return &Config{
		token: getEnv(tgToken, ""),
		port:  getEnv(port, ""),
		debug: getEnvAsBool(debug, false),
	}
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
