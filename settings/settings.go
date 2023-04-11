package settings

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Server connection
var ENVIRONTMENT string
var SERVER_HOST string
var SERVER_PORT string
var GIN_MODE string

// Postgresql connection
var POSTGRESQL_HOST string
var POSTGRESQL_USER string
var POSTGRESQL_PASSWORD string
var POSTGRESQL_DATABASE string
var POSTGRESQL_PORT string
var POSTGRESQL_SSL_MODE string
var POSTGRESQL_TIMEZONE string

// JWT secret
var JWT_SECRET string
var ACCESS_TOKEN_EXPIRE_MINUTES int

func EnvToInt(key string) (int, error) {
	valueString := os.Getenv(key)
	valueInt, err := strconv.Atoi(valueString)
	return valueInt, err
}

func InitiateSettings(pathToEnvFile string) {
	var err error
	if os.Getenv("ENVIRONTMENT") != "PROD" {
		err = godotenv.Load(pathToEnvFile)
		if err != nil {
			panic(err)
		}
	}

	ENVIRONTMENT = os.Getenv("ENVIRONTMENT")
	SERVER_HOST = os.Getenv("SERVER_HOST")
	SERVER_PORT = os.Getenv("SERVER_PORT")
	GIN_MODE = os.Getenv("GIN_MODE")
	POSTGRESQL_HOST = os.Getenv("POSTGRESQL_HOST")
	POSTGRESQL_USER = os.Getenv("POSTGRESQL_USER")
	POSTGRESQL_PASSWORD = os.Getenv("POSTGRESQL_PASSWORD")
	POSTGRESQL_DATABASE = os.Getenv("POSTGRESQL_DATABASE")
	POSTGRESQL_PORT = os.Getenv("POSTGRESQL_PORT")
	POSTGRESQL_SSL_MODE = os.Getenv("POSTGRESQL_SSL_MODE")
	POSTGRESQL_TIMEZONE = os.Getenv("POSTGRESQL_TIMEZONE")
	JWT_SECRET = os.Getenv("JWT_SECRET")
	ACCESS_TOKEN_EXPIRE_MINUTES, err = EnvToInt("ACCESS_TOKEN_EXPIRE_MINUTES")
	if err != nil {
		panic("ACCESS_TOKEN_EXPIRE_MINUTES not defined on env or not a number")
	}
}
