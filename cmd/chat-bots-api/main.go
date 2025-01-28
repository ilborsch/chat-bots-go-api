package main

import (
	_ "chat-bots-api/cmd/chat-bots-api/docs"
	"chat-bots-api/internal/config"
	"chat-bots-api/internal/controller"
	"chat-bots-api/internal/logger"
	"chat-bots-api/internal/repository/mysql"
	"chat-bots-api/internal/repository/sso"
	"chat-bots-api/internal/router"
	"chat-bots-api/internal/usecase"
)

// @title Bot Factory API
// @version 1.0
// @description Bot Factory backend API written in Go Gin
// @securityDefinitions.apikey BearerAuth
// @in cookie
// @name Authorization
// @host localhost:8083
// @BasePath /api/v1/
func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("starting application")

	// Repository is responsible for data layer in the Clean Architecture
	mySQL := mysql.New(cfg.MySQLConfig.Username, cfg.MySQLConfig.Password, cfg.MySQLConfig.Host, cfg.MySQLConfig.Port)
	userRepo := mysql.NewUserRepository(mySQL)
	chatBotRepo := mysql.NewChatBotRepository(mySQL)
	fileRepo := mysql.NewFileRepository(mySQL)
	ssoRepo := sso.NewRepository(cfg.SSOConfig.Host, cfg.SSOConfig.Port)

	// Usecase is responsible for business layer in the Clean Architecture
	useCase := usecase.New(log, cfg.OpenAIConfig.APIKey, userRepo, chatBotRepo, fileRepo, ssoRepo)

	// Controller is responsible for application input/output in the Clean Architecture
	control := controller.New(log, useCase)

	r := router.New(
		log,
		cfg.SSOConfig.Secret,
		control,
		cfg.SSOConfig.AppID,
		cfg.SSOConfig.Timeout,
	)

	const host = "0.0.0.0"
	r.Run(host, cfg.Port)
}
