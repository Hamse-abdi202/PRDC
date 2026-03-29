package infra

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port       string
	DBHost     string
	DBUser     string
	DBPassword string
	DBPort     string
	DBName     string
}

var Configurations AppConfig

func InitEnv() {
	err := godotenv.Load()
	godotenv.Load()
	godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load env variables")
	}

	Configurations.Port = os.Getenv("PORT")
	Configurations.DBHost = os.Getenv("DB_HOST")
	Configurations.DBUser = os.Getenv("DB_USER")
	Configurations.DBPassword = os.Getenv("DB_PASSWORD")
	Configurations.DBPort = os.Getenv("DB_PORT")
	Configurations.DBName = os.Getenv("DB_NAME")
}
