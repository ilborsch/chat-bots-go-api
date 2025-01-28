package router

import (
	_ "chat-bots-api/cmd/chat-bots-api/docs"
	"chat-bots-api/internal/controller"
	"chat-bots-api/internal/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"time"
)

type Router struct {
	log     *slog.Logger
	engine  *gin.Engine
	control *controller.Controller
}

func New(
	log *slog.Logger,
	jwtSecret string,
	control *controller.Controller,
	appID int,
	tokenTTL time.Duration,
) *Router {
	r := gin.New()
	setupMiddleware(r)
	setupRoutes(r, control, log, jwtSecret, appID, tokenTTL)
	return &Router{
		log:     log,
		engine:  r,
		control: control,
	}
}

func (r *Router) Run(host string, port int) {
	addr := fmt.Sprintf("%s:%v", host, port)
	if err := r.engine.Run(addr); err != nil {
		panic("error running application: " + err.Error())
	}
}

func setupMiddleware(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
}

func setupRoutes(
	r *gin.Engine,
	control *controller.Controller,
	log *slog.Logger,
	jwtSecret string,
	appID int,
	tokenTTL time.Duration,
) {

	apiRouter := r.Group("/api/v1")
	setupGeneralRoutes(apiRouter, control, appID, tokenTTL)
	setupUserRoutes(apiRouter, control, log, jwtSecret)
	setupChatBotRoutes(apiRouter, control, log, jwtSecret)
	setupFileRoutes(apiRouter, control, log, jwtSecret)

	//swagger
	r.GET("/docs/*any", gin.WrapH(httpSwagger.WrapHandler))
}

func setupGeneralRoutes(r *gin.RouterGroup, c *controller.Controller, appID int, tokenTTL time.Duration) {
	r.POST("/login", c.SSOController.Login(appID, tokenTTL))
	r.POST("/register", c.SSOController.Register())
}

func setupUserRoutes(r *gin.RouterGroup, c *controller.Controller, log *slog.Logger, jwtSecret string) {
	userGroup := r.Group("/user", middleware.WithJWTAuth(log, jwtSecret))

	userGroup.GET("/", c.UserController.GetUser())
	userGroup.PUT("/", c.UserController.UpdatePlan())
	userGroup.GET("/chat-bots", c.ChatBotController.GetUserChatBots())
}

func setupChatBotRoutes(r *gin.RouterGroup, c *controller.Controller, log *slog.Logger, jwtSecret string) {
	chatBotGroup := r.Group("/chat-bot", middleware.WithJWTAuth(log, jwtSecret))

	chatBotGroup.GET("/:id", c.ChatBotController.GetByID())
	chatBotGroup.POST("/", c.ChatBotController.SaveChatBot())
	chatBotGroup.PATCH("/:id", c.ChatBotController.UpdateChatBot())
	chatBotGroup.DELETE("/:id", c.ChatBotController.RemoveChatBot())
	chatBotGroup.GET("/:id/chat/ws", c.ChatBotController.StartChat())

	chatBotGroup.GET("/:id/files", c.FileController.GetChatBotFiles())
	chatBotGroup.POST("/:id/file", c.FileController.SaveFile())
}

func setupFileRoutes(r *gin.RouterGroup, c *controller.Controller, log *slog.Logger, jwtSecret string) {
	fileGroup := r.Group("/file", middleware.WithJWTAuth(log, jwtSecret))

	fileGroup.GET("/:id", c.FileController.GetByID())
	fileGroup.DELETE("/:id", c.FileController.RemoveFile())
}
