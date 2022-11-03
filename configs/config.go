package configs

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBDriver string

	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBNameTenant string

	RsaPrivateKey string
	RsaPublicKey  string

	RedisPort     string
	RedisHostname string
	RedisPassword string
	RedisDatabase string

	EmailHost     string
	EmailUsername string
	EmailPassword string
	EmailPort     string

	ElasticUrls     []string
	ElasticUser     string
	ElasticPassword string
}

var config *Config

func Init() {
	err := godotenv.Load("./configs/.env")
	if err != nil {
		log.Printf("[Init] error on loading env from file: %+v", err)
	}
	var elasticUrls string

	elasticUrls = os.Getenv("ELASTIC_URLS")

	config = &Config{
		AppPort: os.Getenv("PORT"),

		DBDriver:     os.Getenv("DB_DRIVER"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_DBNAME"),
		DBNameTenant: os.Getenv("DB_DBNAME_TENANT"),

		RsaPrivateKey: os.Getenv("RSA_PRIVATE_KEY"),
		RsaPublicKey:  os.Getenv("RSA_PUBLIC_KEY"),

		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisHostname: os.Getenv("REDIS_HOSTNAME"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDatabase: os.Getenv("REDIS_DATABASE"),

		EmailHost:     os.Getenv("EMAIL_HOST"),
		EmailUsername: os.Getenv("EMAIL_USERNAME"),
		EmailPassword: os.Getenv("EMAIL_PASSWORD"),
		EmailPort:     os.Getenv("EMAIL_PORT"),

		ElasticUrls:     strings.Split(elasticUrls, ","),
		ElasticUser:     os.Getenv("ELASTIC_USER"),
		ElasticPassword: os.Getenv("ELASTIC_PASSWORD"),
	}

	// if config.AppName == "" {
	// 	log.Panicf("[Init] app name cannot be empty")
	// }

	// if config.Port == "" {
	// 	log.Panicf("[Init] app address cannot be empty")
	// }

	if config.DBDriver == "" {
		log.Panicf("[Init] db driver cannot be empty")
	}

	if config.DBDriver == "sqlite3" {
		if config.DBName == "" || config.DBNameTenant == "" {
			log.Panicf("[Init] db name  db name tenant cannot be empty")
		}
	} else {
		if config.DBDriver == "" || config.DBHost == "" || config.DBPort == "" || config.DBName == "" || config.DBNameTenant == "" {
			log.Panicf("[Init] db name or address or db name tenant cannot be empty")
		}
	}

}

func Get() *Config {
	return config
}
