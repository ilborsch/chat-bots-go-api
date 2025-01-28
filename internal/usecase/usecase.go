package usecase

import (
	"chat-bots-api/internal/repository"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/ilborsch/openai-go/openai"
	"log/slog"
)

//go:generate mockgen -destination=../../mocks/mock_openai.go -package=mocks github.com/ilborsch/openai-go/openai OpenAIClient

type Usecase struct {
	openAIClient openai.OpenAIClient
	log          *slog.Logger
	repository.UserRepository
	repository.ChatBotRepository
	repository.FileRepository
	repository.SSORepository
}

func New(
	log *slog.Logger,
	openaiAPIKey string,
	userRepo repository.UserRepository,
	cbRepo repository.ChatBotRepository,
	fileRepo repository.FileRepository,
	ssoRepo repository.SSORepository,
) *Usecase {
	return &Usecase{
		log:               log,
		openAIClient:      openai.New(openaiAPIKey),
		UserRepository:    userRepo,
		ChatBotRepository: cbRepo,
		FileRepository:    fileRepo,
		SSORepository:     ssoRepo,
	}
}
