package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	serverPort = "SERVER_PORT"

	dbHost     = "QA_DB_HOST"
	dbPort     = "QA_DB_PORT"
	dbUser     = "QA_DB_USER"
	dbPassword = "QA_DB_PASSWORD"
	dbName     = "QA_DB_NAME"
)

var singleInstance *Config

type Config struct {
	serverPort string

	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
}

func (c Config) ServerPort() string {
	return c.serverPort
}

func (c Config) DbHost() string {
	return c.dbHost
}

func (c Config) DbPort() string {
	return c.dbPort
}

func (c Config) DbUser() string {
	return c.dbUser
}

func (c Config) DbPassword() string {
	return c.dbPassword
}

func (c Config) DbName() string {
	return c.dbName
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("No .env file found")
	}

	singleInstance = new()

	validateConfig()
}

func GetConfig() *Config {
	return singleInstance
}

func new() *Config {
	return &Config{
		serverPort: getEnv(serverPort, ""),

		dbHost:     getEnv(dbHost, "localhost"),
		dbPort:     getEnv(dbPort, "5432"),
		dbUser:     getEnv(dbUser, ""),
		dbPassword: getEnv(dbPassword, ""),
		dbName:     getEnv(dbName, ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}

func validateConfig() {
	c := GetConfig()
	if len(c.ServerPort()) == 0 {
		log.Fatal("Config error: server port is empty")
	}
	if len(c.DbHost()) == 0 {
		log.Fatal("Config error: qa db host is empty")
	}
	if len(c.DbPort()) == 0 {
		log.Fatal("Config error: qa db port is empty")
	}
	if len(c.DbUser()) == 0 {
		log.Fatal("Config error: qa db user is empty")
	}
	if len(c.DbPassword()) == 0 {
		log.Fatal("Config error: qa db password is empty")
	}
	if len(c.DbName()) == 0 {
		log.Fatal("Config error: qa db name is empty")
	}
}
