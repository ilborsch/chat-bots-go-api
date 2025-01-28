package controller

import (
	"chat-bots-api/internal/schemas"
	"chat-bots-api/internal/usecase"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"strconv"
)

type ChatBotController struct {
	log *slog.Logger
	usecase.ChatBotUsecase
}

// GetByID godoc
// @Summary Fetches chat-bot by ID
// @Description Fetch chat-bot object from database by its ID
// @Param id path int true "Chat-bot ID"
// @Produce application/json
// @Security BearerAuth
// @Tags chat-bot
// @Success 200 {object} schemas.ChatBot
// @Router /api/v1/chat-bot/{id} [get]
func (con *ChatBotController) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			con.log.Error("error parsing id")
			respondWithError(c, http.StatusBadRequest, "Invalid ID provided")
			return
		}

		uid, ok := c.Get("uid")
		if !ok {
			con.log.Error("couldn't get userID from jwt token in handler")
			respondWithError(c, http.StatusUnauthorized, "Permission denied")
			return
		}
		userID := int64(uid.(float64))

		bot, err := con.ChatBotUsecase.ChatBot(c, id, userID)
		if err != nil {
			con.log.Error(err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		chatBotSchema := schemas.ChatBot{
			ID:           bot.ID,
			Name:         bot.Name,
			Description:  bot.Description,
			Instructions: bot.Instructions,
		}

		c.JSON(http.StatusOK, chatBotSchema)
	}
}

// SaveChatBot godoc
// @Summary Saves chat-bot
// @Description Save chat-bot object to database
// @Param chat_bot body schemas.SaveChatBotRequest true "Chat bot"
// @Produce application/json
// @Security BearerAuth
// @Tags chat-bot
// @Success 200 {object} schemas.SaveChatBotResponse
// @Router /api/v1/chat-bot/ [post]
func (con *ChatBotController) SaveChatBot() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request schemas.SaveChatBotRequest
		if err := c.BindJSON(&request); err != nil {
			con.log.Error("error parsing request data: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Bad request")
			return
		}

		uid, ok := c.Get("uid")
		if !ok {
			con.log.Error("couldn't get userID from jwt token in handler")
			respondWithError(c, http.StatusUnauthorized, "Permission denied")
			return
		}
		userID := int64(uid.(float64))

		id, err := con.ChatBotUsecase.SaveChatBot(
			c,
			request.Name,
			request.Description,
			request.Instructions,
			userID,
		)
		if err != nil {
			con.log.Error(err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		response := schemas.SaveChatBotResponse{
			ChatBotID: id,
		}
		c.JSON(http.StatusOK, response)
	}
}

// UpdateChatBot godoc
// @Summary Updates chat-bot
// @Description Update chat-bot object in database
// @Param chat_bot body schemas.UpdateChatBotRequest true "Chat bot"
// @Param id path int true "Chat bot"
// @Produce application/json
// @Security BearerAuth
// @Tags chat-bot
// @Success 200 {object} schemas.UpdateChatBotResponse
// @Router /api/v1/chat-bot/{id} [patch]
func (con *ChatBotController) UpdateChatBot() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			con.log.Error(err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		var request schemas.UpdateChatBotRequest
		if err := c.BindJSON(&request); err != nil {
			con.log.Error("error parsing request data: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Bad request")
			return
		}

		uid, ok := c.Get("uid")
		if !ok {
			con.log.Error("couldn't get userID from jwt token in handler")
			respondWithError(c, http.StatusUnauthorized, "Permission denied")
			return
		}
		userID := int64(uid.(float64))

		if err := con.ChatBotUsecase.UpdateChatBot(
			c,
			id,
			userID,
			request.Name,
			request.Description,
			request.Instructions,
		); err != nil {
			con.log.Error(err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		response := schemas.UpdateChatBotResponse{
			Success: true,
		}
		c.JSON(http.StatusOK, response)
	}
}

// RemoveChatBot godoc
// @Summary Removes chat-bot
// @Description Removes chat-bot object from database
// @Param id path int true "Chat bot"
// @Produce application/json
// @Security BearerAuth
// @Tags chat-bot
// @Success 200 {object} schemas.RemoveChatBotResponse
// @Router /api/v1/chat-bot/{id} [delete]
func (con *ChatBotController) RemoveChatBot() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			con.log.Error("error parsing id " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid ID provided")
			return
		}

		uid, ok := c.Get("uid")
		if !ok {
			con.log.Error("couldn't get userID from jwt token in handler")
			respondWithError(c, http.StatusUnauthorized, "Permission denied")
			return
		}
		userID := int64(uid.(float64))

		if err := con.ChatBotUsecase.RemoveChatBot(c, id, userID); err != nil {
			con.log.Error(err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		response := schemas.RemoveChatBotResponse{
			Success: true,
		}
		c.JSON(http.StatusOK, response)
	}
}

// StartChat godoc
// @Summary Start chat with chat-bot
// @Description Starts chat with a chat-bot with id = chatBotID, opens a websocket connection
// @Param id path int true "Chat bot"
// @Security BearerAuth
// @Tags chat-bot
// @Success 101 {string} string "Switching Protocols"
// @Router /api/v1/chat-bot/{id}/chat/ws [get]
func (con *ChatBotController) StartChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		chatBotID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			con.log.Error("error parsing id " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid ID provided")
			return
		}

		uid, ok := c.Get("uid")
		if !ok {
			con.log.Error("couldn't get userID from jwt token in handler")
			respondWithError(c, http.StatusUnauthorized, "Permission denied")
			return
		}
		userID := int64(uid.(float64))

		threadID, err := con.ChatBotUsecase.StartChat(c, chatBotID, userID)
		if err != nil {
			con.log.Error(err.Error())
			respondWithError(c, http.StatusBadRequest, "Bad request")
			return
		}

		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			con.log.Error("error upgrading request to web-socket: " + err.Error())
			respondWithError(c, http.StatusInternalServerError, "Failed to upgrade to WebSocket")
			return
		}
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					con.log.Info("websocket connection closed")
					return
				}
				con.log.Error("Error reading message: " + err.Error())
				return
			}

			var request schemas.SendMessageRequest
			fmt.Println(string(msg))
			if err := json.Unmarshal(msg, &request); err != nil {
				sendErrorWSMessage(conn, "Invalid request format")
				return
			}

			responseChan := make(chan usecase.AssistantResponse)
			go con.ChatBotUsecase.SendMessage(
				c,
				chatBotID,
				userID,
				threadID,
				request.Message,
				responseChan,
			)

			assistantResponse := <-responseChan
			if assistantResponse.Err != nil {
				sendErrorWSMessage(conn, "Could not retrieve chat-bot response")
				return
			}

			response := schemas.SendMessageResponse{
				Response: assistantResponse.Content,
			}
			respJSON, err := json.Marshal(response)
			if err != nil {
				sendErrorWSMessage(conn, "Could not retrieve chat-bot response")
				return
			}

			conn.WriteMessage(websocket.TextMessage, respJSON)
		}
	}
}

// GetUserChatBots godoc
// @Summary Get all user chat-bots
// @Description Returns all user chat-bots from database
// @Security BearerAuth
// @Tags user
// @Success 200 {object} schemas.UserChatBots
// @Router /api/v1/user/chat-bots [get]
func (con *ChatBotController) GetUserChatBots() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get("uid")
		if !ok {
			con.log.Error("couldn't get userID from jwt token in handler")
			respondWithError(c, http.StatusUnauthorized, "Permission denied")
			return
		}
		userID := int64(uid.(float64))

		chatBots, err := con.ChatBotUsecase.UserChatBots(c, userID)
		if err != nil {
			con.log.Error("error retrieving chat-bots: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		chatBotsPayload := make([]schemas.ChatBot, len(chatBots))
		for i, chatBot := range chatBots {
			chatBotsPayload[i] = schemas.ChatBot{
				ID:           chatBot.ID,
				Name:         chatBot.Name,
				Description:  chatBot.Description,
				Instructions: chatBot.Instructions,
			}
		}
		payload := schemas.UserChatBots{
			ChatBots: chatBotsPayload,
		}
		c.JSON(http.StatusOK, payload)
	}
}
