package controller

import (
	"chat-bots-api/domain"
	"chat-bots-api/internal/schemas"
	"chat-bots-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type SSOController struct {
	log *slog.Logger
	usecase.SSOUsecase
	usecase.UserUsecase
}

// Register godoc
// @Summary User registration
// @Description Register user with email and password
// @Param credentials body schemas.RegisterRequest true "Register credentials"
// @Tags auth
// @Success 200 {object} schemas.RegisterResponse
// @Router /api/v1/register/ [post]
func (con *SSOController) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request schemas.RegisterRequest
		if err := c.BindJSON(&request); err != nil {
			con.log.Error("error parsing request data: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Bad request")
			return
		}

		uid, err := con.UserUsecase.SaveUser(c, request.Email, request.Password, domain.FreePlan)
		if err != nil {
			con.log.Error("error saving new user: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Bad request")
			return
		}

		if _, err := con.SSOUsecase.Register(c, request.Email, request.Password, uid); err != nil {
			con.log.Error("error registering new user: " + err.Error())
			respondWithError(c, http.StatusInternalServerError, "Internal server error")
			return
		}

		response := schemas.RegisterResponse{
			ID: uid,
		}
		c.JSON(http.StatusOK, response)
	}
}

// Login godoc
// @Summary User login
// @Description Login user with email and password
// @Param credentials body schemas.LoginRequest true "Login credentials"
// @Tags auth
// @Success 200 {object} schemas.LoginResponse
// @Router /api/v1/login/ [post]
func (con *SSOController) Login(appID int, tokenTTL time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request schemas.LoginRequest
		if err := c.BindJSON(&request); err != nil {
			con.log.Error("error parsing request data: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Bad request")
			return
		}

		token, err := con.SSOUsecase.Login(c, request.Email, request.Password, appID)
		if err != nil {
			con.log.Error("error logging user in: " + err.Error())
			respondWithError(c, http.StatusBadRequest, "Bad request")
			return
		}

		response := schemas.LoginResponse{
			Token: token,
		}
		c.SetCookie(
			"Authorization",
			token,
			int(tokenTTL.Seconds()),
			"/",
			"localhost",
			false,
			true,
		)
		c.JSON(http.StatusOK, response)
	}
}
