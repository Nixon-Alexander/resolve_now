package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Nixon-Alexander/resolve_now.git/config"
	"github.com/Nixon-Alexander/resolve_now.git/controller"
	"github.com/Nixon-Alexander/resolve_now.git/repository"
	"github.com/Nixon-Alexander/resolve_now.git/router"
	"github.com/joho/godotenv"
)

func main() {
	logger := config.InitLog()

	// Membaca .env dari command di make file
	envFile := flag.String("env", ".env", "env file")
	flag.Parse()

	err := godotenv.Load(*envFile)
	if err != nil {
		logger.Error("Error loading .env file")
	}

	// Start QDrant Vector DB
	qdrants, err := config.InitQDrant()
	if err != nil {
		logger.Error("Error loading qdrant", "error", err)
		os.Exit(1)
	}

	if err := qdrants.CreateChunksCollection(); err != nil {
		logger.Error("Error creating qdrant chunks collection", "error", err)
		return
	}
	logger.Info("QDrant is connected")

	// Injection
	db, err := config.InitMySqlDatabase(logger).Connect()
	if err != nil {
		logger.Error("Database connection failed", "error", err)
		return
	}

	err = db.UpDB()
	if err != nil {
		logger.Error("Error running migration", "error", err)
		return
	}

	companyRepo := repository.InitCompanyRepository(db, logger)
	companyController := controller.InitCompanyController(companyRepo, logger)

	documentRepo := repository.InitDocumentRepository(db, logger, qdrants)
	documentContoller := controller.InitDocumentController(documentRepo, logger)

	chatRepo := repository.InitChatRepository(db, qdrants, logger)
	chatController := controller.InitChatController(chatRepo, logger)

	r := router.InitRoutes(companyController, documentContoller, chatController)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	var server = fmt.Sprint(host, ":", port)
	logger.Info("Server Started", "host", host, "port", port)

	r.Run(server)
}
