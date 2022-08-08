package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	tgToken                     = "TELEGRAM_BOT_API_TOKEN"
	serverPort                  = "SERVER_PORT"
	restPort                    = "REST_PORT"
	restGrpcPort                = "REST_GRPC_PORT"
	debug                       = "DEBUG"
	requestTimeOutInMilliSecond = "REQUEST_TIMEOUT_IN_MILLISECOND"

	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbUser     = "DB_USER"
	dbPassword = "DB_PASSWORD"
	dbName     = "DB_NAME"
)

var singleInstance *Config

type Config struct {
	token                       string
	serverPort                  string
	restPort                    string
	restGrpcPort                string
	debug                       bool
	requestTimeOutInMilliSecond time.Duration

	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
}

func (c Config) Token() string {
	return c.token
}

func (c Config) ServerPort() string {
	return c.serverPort
}

func (c Config) RestPort() string {
	return c.restPort
}

func (c Config) RestGrpcPort() string {
	return c.restGrpcPort
}

func (c Config) Debug() bool {
	return c.debug
}

func (c Config) RequestTimeOutInMilliSecond() time.Duration {
	return c.requestTimeOutInMilliSecond
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
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	singleInstance = new()

	validateConfig()
}

func GetConfig() *Config {
	return singleInstance
}

func new() *Config {
	requestTimeOutInMilliSecond := getEnvAsInt64(requestTimeOutInMilliSecond, 500)
	return &Config{
		token:        getEnv(tgToken, ""),
		serverPort:   getEnv(serverPort, ""),
		restPort:     getEnv(restPort, ""),
		restGrpcPort: getEnv(restGrpcPort, ""),

		dbHost:     getEnv(dbHost, "localhost"),
		dbPort:     getEnv(dbPort, "5432"),
		dbUser:     getEnv(dbUser, ""),
		dbPassword: getEnv(dbPassword, ""),
		dbName:     getEnv(dbName, ""),

		debug:                       getEnvAsBool(debug, false),
		requestTimeOutInMilliSecond: time.Duration(requestTimeOutInMilliSecond) * time.Millisecond,
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

func getEnvAsInt64(name string, defaultValue int64) int64 {
	value := getEnv(name, "")
	if valueInt, err := strconv.Atoi(value); err == nil {
		return int64(valueInt)
	}
	return defaultValue
}

func validateConfig() {
	c := GetConfig()
	if len(c.ServerPort()) == 0 {
		log.Fatal("Config error: server port is empty")
	}
	if len(c.RestPort()) == 0 {
		log.Fatal("Config error: rest port is empty")
	}
	if len(c.RestGrpcPort()) == 0 {
		log.Fatal("Config error: rest grpc port is empty")
	}
	if len(c.DbHost()) == 0 {
		log.Fatal("Config error: db host is empty")
	}
	if len(c.DbPort()) == 0 {
		log.Fatal("Config error: db port is empty")
	}
	if len(c.DbUser()) == 0 {
		log.Fatal("Config error: db user is empty")
	}
	if len(c.DbPassword()) == 0 {
		log.Fatal("Config error: db password is empty")
	}
	if len(c.DbName()) == 0 {
		log.Fatal("Config error: db name is empty")
	}
}
