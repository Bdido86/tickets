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
	clientPort                  = "CLIENT_PORT"
	clientGrpcPort              = "CLIENT_GRPC_PORT"
	debugLevel                  = "DEBUG_LEVEL"
	requestTimeOutInMilliSecond = "REQUEST_TIMEOUT_IN_MILLISECOND"

	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbUser     = "DB_USER"
	dbPassword = "DB_PASSWORD"
	dbName     = "DB_NAME"

	redisAddr     = "REDIS_ADDR"
	redisPassword = "REDIS_PASSWORD"
	redisDb       = "REDIS_DB"
)

var singleInstance *Config

type Config struct {
	token                       string
	serverPort                  string
	clientPort                  string
	clientGrpcPort              string
	debugLevel                  string
	requestTimeOutInMilliSecond time.Duration

	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string

	redisAddr     string
	redisPassword string
	redisDb       int
}

func (c Config) Token() string {
	return c.token
}

func (c Config) ServerPort() string {
	return c.serverPort
}

func (c Config) ClientPort() string {
	return c.clientPort
}

func (c Config) ClientGrpcPort() string {
	return c.clientGrpcPort
}

func (c Config) DebugLevel() string {
	return c.debugLevel
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

func (c Config) RedisAddr() string {
	return c.redisAddr
}

func (c Config) RedisPassword() string {
	return c.redisPassword
}

func (c Config) RedisDb() int {
	return c.redisDb
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
		token:          getEnv(tgToken, ""),
		serverPort:     getEnv(serverPort, ""),
		clientPort:     getEnv(clientPort, ""),
		clientGrpcPort: getEnv(clientGrpcPort, ""),

		dbHost:     getEnv(dbHost, "localhost"),
		dbPort:     getEnv(dbPort, "5432"),
		dbUser:     getEnv(dbUser, ""),
		dbPassword: getEnv(dbPassword, ""),
		dbName:     getEnv(dbName, ""),

		redisAddr:     getEnv(redisAddr, ""),
		redisPassword: getEnv(redisPassword, ""),
		redisDb:       getEnvAsInt(redisDb, 0),

		debugLevel:                  getEnv(debugLevel, "info"),
		requestTimeOutInMilliSecond: time.Duration(requestTimeOutInMilliSecond) * time.Millisecond,
	}
}

func getEnv(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultValue int) int {
	value := getEnv(name, "")
	if valueInt, err := strconv.Atoi(value); err == nil {
		return valueInt
	}
	return defaultValue
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
	if len(c.ClientPort()) == 0 {
		log.Fatal("Config error: client port is empty")
	}
	if len(c.ClientGrpcPort()) == 0 {
		log.Fatal("Config error: client grpc port is empty")
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
	if len(c.RedisAddr()) == 0 {
		log.Fatal("Config error: redis addr is empty")
	}
}
