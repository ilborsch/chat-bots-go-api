package repository

import (
	domain "chat-bots-api/domain"
	"context"
	"github.com/ilborsch/sso-proto/gen/go/sso"
)

//go:generate mockgen -destination=../../mocks/mock_chat_bot_repository.go -package=mocks chat-bots-api/internal/repository ChatBotRepository
//go:generate mockgen -destination=../../mocks/mock_file_repository.go -package=mocks chat-bots-api/internal/repository FileRepository
//go:generate mockgen -destination=../../mocks/mock_user_repository.go -package=mocks chat-bots-api/internal/repository UserRepository
//go:generate mockgen -destination=../../mocks/mock_sso_repository.go -package=mocks chat-bots-api/internal/repository SSORepository

type ChatBotRepository interface {
	ChatBot(ctx context.Context, id, ownerID int64) (domain.ChatBot, error)
	UserChatBots(ctx context.Context, ownerID int64) ([]domain.ChatBot, error)
	SaveChatBot(ctx context.Context, chatBot domain.ChatBot) (id int64, err error)
	UpdateChatBot(ctx context.Context, id, ownerID int64, newChatBot domain.ChatBot) error
	RemoveChatBot(ctx context.Context, id int64, userID int64) error
}

type FileRepository interface {
	File(ctx context.Context, id, ownerID int64) (domain.File, error)
	ChatBotFiles(ctx context.Context, chatBotID, ownerID int64) ([]domain.File, error)
	SaveFile(ctx context.Context, file domain.File, ownerID int64) (id int64, err error)
	RemoveFile(ctx context.Context, id int64, fileSize int, ownerID int64) error
}

type UserRepository interface {
	User(ctx context.Context, id int64) (domain.User, error)
	UserByEmail(ctx context.Context, email string) (domain.User, error)
	SaveUser(ctx context.Context, user domain.User) (id int64, err error)
	UpdatePlan(ctx context.Context, id int64, newUser domain.User) error
	UpdateMessagesLeft(ctx context.Context, id int64) error
}

type SSORepository interface {
	IsAdmin(ctx context.Context, request *sso.IsAdminRequest) (bool, error)
	Login(ctx context.Context, request *sso.LoginRequest) (string, error)
	Register(ctx context.Context, request *sso.RegisterRequest) (int64, error)
}
