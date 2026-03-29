package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/ilwaad12/product-project/infra"
	"github.com/ilwaad12/product-project/routes"
	"github.com/joho/godotenv"
)

func main() {
	slog.Info("Initialliaing env variables")
	infra.InitEnv()
	config := infra.Configurations

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slog.Info("Connecting to database")
	infra.ConnectDB()
	slog.Info("✅ Connected to database successfully")

	r := gin.Default()

	routes.RegisterRoutes(r)

	slog.Info("✅ Application is runing on port :5000")
	r.Run(fmt.Sprintf(":%s", config.Port))
}
