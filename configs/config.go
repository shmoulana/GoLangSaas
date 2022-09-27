package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	DBDriver      string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBNameTenant  string
	RsaPrivateKey string
	RsaPublicKey  string
}

var config *Config

func Init() {
	err := godotenv.Load("./configs/.env")
	if err != nil {
		log.Printf("[Init] error on loading env from file: %+v", err)
	}

	config = &Config{
		AppPort:       os.Getenv("APP_PORT"),
		DBDriver:      os.Getenv("DB_DRIVER"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_DBNAME"),
		DBNameTenant:  os.Getenv("DB_DBNAME_TENANT"),
		RsaPrivateKey: os.Getenv("RSA_PRIVATE_KEY"),
		RsaPublicKey:  os.Getenv("RSA_PUBLIC_KEY"),
	}

	// if config.AppName == "" {
	// 	log.Panicf("[Init] app name cannot be empty")
	// }

	// if config.Port == "" {
	// 	log.Panicf("[Init] app address cannot be empty")
	// }

	if config.DBDriver == "" || config.DBHost == "" || config.DBPort == "" || config.DBName == "" || config.DBNameTenant == "" {
		log.Panicf("[Init] db name or address or db driver or db name tenant cannot be empty")
	}
}

func Get() *Config {
	return config
}
