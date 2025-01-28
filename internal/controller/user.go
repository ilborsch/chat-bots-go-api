package controller

import (
	"chat-bots-api/internal/schemas"
	"chat-bots-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type UserController struct {
	log *slog.Logger
	usecase.UserUsecase
}

func (con *UserController) GetByEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")
		if email == "" {
			con.log.Error("empty email provided")
			respondWithError(c, http.StatusOK, "Empty email provided.")
			return
		}

		user, err := con.UserUsecase.UserByEmail(c, email)
		if err != nil {
			con.log.Error("empty email provided")
			respondWithError(c, http.StatusOK, "Empty email provided.")
			return
		}

		payload := schemas.User{
			ID:             user.ID,
			Email:          user.Email,
			Plan:           user.Plan,
			PlanBoughtDate: user.PlanBoughtDate,
			MessagesLeft:   user.MessagesLeft,
			BytesDataLeft:  user.BytesDataLeft,
			BotsLeft:       user.BotsLeft,
		}
		c.JSON(http.StatusOK, payload)
	}
}

// GetUser godoc
// @Summary Fetches user by token
// @Description Fetch user data from database by session token
// @Produce application/json
// @Security BearerAuth
// @Tags user
// @Success 200 {object} schemas.User
// @Router /api/v1/user/ [get]
func (con *UserController) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get("uid")
		if !ok {
			con.log.Error("couldn't get userID from jwt token in handler")
			respondWithError(c, http.StatusUnauthorized, "Permission denied")
			return
		}
		userID := int64(uid.(float64))

		user, err := con.UserUsecase.User(c, userID)
		if err != nil {
			con.log.Error("empty email provided")
			respondWithError(c, http.StatusOK, "Empty email provided.")
			return
		}

		payload := schemas.User{
			ID:             user.ID,
			Email:          user.Email,
			Plan:           user.Plan,
			PlanBoughtDate: user.PlanBoughtDate,
			MessagesLeft:   user.MessagesLeft,
			BytesDataLeft:  user.BytesDataLeft,
			BotsLeft:       user.BotsLeft,
		}
		c.JSON(http.StatusOK, payload)
	}
}

func (con *UserController) SaveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request schemas.SaveUserRequest
		if err := c.BindJSON(&request); err != nil {
			con.log.Error("error parsing request data: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Bad request")
			return
		}

		id, err := con.UserUsecase.SaveUser(c, request.Email, request.Password, request.Plan)
		if err != nil {
			con.log.Error("error saving user: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Bad request")
			return
		}

		payload := schemas.SaveUserResponse{
			UserID: id,
		}
		c.JSON(http.StatusOK, payload)
	}
}

// UpdatePlan godoc
// @Summary Update user plan
// @Description Update user plan (Free: "free_plan", Business: "business_plan", Enterprise: "enterprise_plan")
// @Produce application/json
// @Param new-plan body schemas.UpdatePlanRequest true "Chat-bot ID"
// @Security BearerAuth
// @Tags user
// @Success 200 {object} schemas.UpdatePlanResponse
// @Router /api/v1/user/ [put]
func (con *UserController) UpdatePlan() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get("uid")
		if !ok {
			con.log.Error("couldn't get userID from jwt token in handler")
			respondWithError(c, http.StatusUnauthorized, "Permission denied")
			return
		}
		userID := int64(uid.(float64))

		var request schemas.UpdatePlanRequest
		if err := c.BindJSON(&request); err != nil {
			con.log.Error("error parsing request data: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Bad request")
			return
		}

		if err := con.UserUsecase.UpdatePlan(c, userID, request.Plan); err != nil {
			con.log.Error("error updating plan: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		payload := schemas.UpdatePlanResponse{
			Success: true,
		}
		c.JSON(http.StatusOK, payload)
	}
}
