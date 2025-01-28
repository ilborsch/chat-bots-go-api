package controller

import (
	"chat-bots-api/internal/schemas"
	"chat-bots-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type FileController struct {
	log *slog.Logger
	usecase.FileUsecase
}

// GetByID godoc
// @Summary Returns a file meta-data
// @Description Returns a file meta-data by its ID from the database
// @Param id path int true "File ID"
// @Security BearerAuth
// @Tags file
// @Success 200 {object} schemas.File
// @Router /api/v1/file/{id} [get]
func (con *FileController) GetByID() gin.HandlerFunc {
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

		file, err := con.FileUsecase.File(c, id, userID)
		if err != nil {
			con.log.Error(err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		response := schemas.File{
			ID:        file.ID,
			Filename:  file.Filename,
			ChatBotID: file.ChatBotID,
		}
		c.JSON(http.StatusOK, response)
	}
}

// SaveFile godoc
// @Summary Saves a file
// @Description Saves a file to database and adds it to the chat-bot
// @Accept mpfd
// @Param chat-bot-id path int true "Chat Bot ID"
// @Param file formData file true "File to upload"
// @Security BearerAuth
// @Tags file
// @Success 200 {object} schemas.SaveFileResponse
// @Router /api/v1/chat-bot/{chat-bot-id}/file/ [post]
func (con *FileController) SaveFile() gin.HandlerFunc {
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

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			con.log.Error("error retrieving file data: " + err.Error())
			respondWithError(c, http.StatusInternalServerError, "Internal error")
			return
		}

		filename := header.Filename

		fileData, err := io.ReadAll(file)
		if err != nil {
			con.log.Error("error reading file data: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid file data")
			return
		}

		id, err := con.FileUsecase.SaveFile(
			c,
			chatBotID,
			userID,
			filename,
			fileData,
		)
		if err != nil {
			con.log.Error(err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		response := schemas.SaveFileResponse{
			FileID: id,
		}
		c.JSON(http.StatusOK, response)
	}
}

// RemoveFile godoc
// @Summary Removes a file
// @Description Removes a file from database and chat-bot
// @Param id path int true "File ID"
// @Security BearerAuth
// @Tags file
// @Success 200 {object} schemas.RemoveFileResponse
// @Router /api/v1/file/{id} [delete]
func (con *FileController) RemoveFile() gin.HandlerFunc {
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

		if err := con.FileUsecase.RemoveFile(c, id, userID); err != nil {
			con.log.Error("error removing file: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		response := schemas.RemoveFileResponse{
			Success: true,
		}
		c.JSON(http.StatusOK, response)
	}
}

// GetChatBotFiles godoc
// @Summary Get chat-bot files
// @Description Return chat-bot files meta-data from database
// @Param chat-bot-id path int true "Chat Bot ID"
// @Security BearerAuth
// @Tags file
// @Success 200 {object} schemas.ChatBotFiles
// @Router /api/v1/chat-bot/{chat-bot-id}/files [get]
func (con *FileController) GetChatBotFiles() gin.HandlerFunc {
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

		files, err := con.FileUsecase.ChatBotFiles(c, chatBotID, userID)
		if err != nil {
			con.log.Error("error retrieving chat-bot files: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		filesPayload := make([]schemas.File, len(files))
		for i, file := range files {
			filesPayload[i] = schemas.File{
				ID:        file.ID,
				Filename:  file.Filename,
				ChatBotID: file.ChatBotID,
			}
		}
		payload := schemas.ChatBotFiles{
			Files: filesPayload,
		}
		c.JSON(http.StatusOK, payload)
	}
}
