package controller

import (
	"chat-bots-api/internal/schemas"
	"chat-bots-api/internal/usecase"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log/slog"
)

type Controller struct {
	UserController
	ChatBotController
	FileController
	SSOController
}

func New(log *slog.Logger, useCase *usecase.Usecase) *Controller {
	return &Controller{
		UserController: UserController{
			log:         log,
			UserUsecase: useCase,
		},
		ChatBotController: ChatBotController{
			log:            log,
			ChatBotUsecase: useCase,
		},
		FileController: FileController{
			log:         log,
			FileUsecase: useCase,
		},
		SSOController: SSOController{
			log:         log,
			SSOUsecase:  useCase,
			UserUsecase: useCase,
		},
	}
}

func respondWithError(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, gin.H{"error": errorMessage})
}

func sendErrorWSMessage(conn *websocket.Conn, errorMessage string) {
	payload := schemas.SendMessageError{
		Error: errorMessage,
	}
	errorJSON, _ := json.Marshal(&payload)

	_ = conn.WriteMessage(websocket.TextMessage, errorJSON)
}
