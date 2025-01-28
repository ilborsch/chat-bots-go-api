package usecase

import (
	"chat-bots-api/domain"
	"chat-bots-api/internal/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/ilborsch/openai-go/openai/assistants/runs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestUsecase_ChatBot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBotRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:               logger,
		ChatBotRepository: mockRepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name        string
		id          int64
		ownerID     int64
		mockSetup   func()
		expectError bool
		errorMsg    string
	}{
		{
			name:    "Valid chat-bot",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{ID: 1, OwnerID: 1}, nil)
			},
			expectError: false,
		},
		{
			name:    "Invalid chat-bot ID",
			id:      0,
			ownerID: 1,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectError: true,
			errorMsg:    "invalid chat-bot id provided",
		},
		{
			name:    "Invalid owner ID",
			id:      1,
			ownerID: 0,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectError: true,
			errorMsg:    "invalid owner id provided 0",
		},
		{
			name:    "Repository Error",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{}, errors.New("repository error"))
			},
			expectError: true,
			errorMsg:    "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			chatBot, err := u.ChatBot(ctx, tt.id, tt.ownerID)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, chatBot.ID)
				assert.Equal(t, tt.ownerID, chatBot.OwnerID)
			}
		})
	}
}

func TestUsecase_UserChatBots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBotRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:               logger,
		ChatBotRepository: mockRepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name         string
		ownerID      int64
		mockSetup    func()
		expectError  bool
		errorMsg     string
		expectedBots []domain.ChatBot
	}{
		{
			name:    "Valid user chat-bots",
			ownerID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().UserChatBots(ctx, int64(1)).Return([]domain.ChatBot{{ID: 1, OwnerID: 1}}, nil)
			},
			expectError:  false,
			expectedBots: []domain.ChatBot{{ID: 1, OwnerID: 1}},
		},
		{
			name:    "Invalid owner ID",
			ownerID: 0,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectError:  true,
			errorMsg:     "invalid owner id provided 0",
			expectedBots: nil,
		},
		{
			name:    "Repository error",
			ownerID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().UserChatBots(ctx, int64(1)).Return(nil, errors.New("repository error"))
			},
			expectError:  true,
			errorMsg:     "repository error",
			expectedBots: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			chatBots, err := u.UserChatBots(ctx, tt.ownerID)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, chatBots)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBots, chatBots)
			}
		})
	}
}

func TestUsecase_SaveChatBot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBotRepository(ctrl)
	mockOpenAIClient := mocks.NewMockOpenAIClient(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:               logger,
		ChatBotRepository: mockRepo,
		openAIClient:      mockOpenAIClient,
	}

	ctx := context.TODO()

	tests := []struct {
		name         string
		nameArg      string
		description  string
		instructions string
		ownerID      int64
		mockSetup    func()
		expectError  bool
		errorMsg     string
		expectedID   int64
	}{
		{
			name:         "Valid chat bot creation",
			nameArg:      "TestBot",
			description:  "A test bot",
			instructions: "Test instructions",
			ownerID:      1,
			mockSetup: func() {
				mockOpenAIClient.EXPECT().CreateVectorStore("1 - TestBot").Return("vsID", nil)
				mockOpenAIClient.EXPECT().CreateAssistant("1 - TestBot", "Test instructions", "vsID", nil).Return("assistantID", nil)
				mockRepo.EXPECT().SaveChatBot(ctx, gomock.Any()).Return(int64(1), nil)
			},
			expectError: false,
			expectedID:  1,
		},
		{
			name:         "Validation error",
			nameArg:      "",
			description:  "A test bot",
			instructions: "Test instructions",
			ownerID:      1,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectError: true,
			errorMsg:    "empty chat-bot name provided",
			expectedID:  0,
		},
		{
			name:         "Create VectorStore error",
			nameArg:      "TestBot",
			description:  "A test bot",
			instructions: "Test instructions",
			ownerID:      1,
			mockSetup: func() {
				mockOpenAIClient.EXPECT().CreateVectorStore("1 - TestBot").Return("", errors.New("vector store error"))
			},
			expectError: true,
			errorMsg:    "vector store error",
			expectedID:  0,
		},
		{
			name:         "Create Assistant error",
			nameArg:      "TestBot",
			description:  "A test bot",
			instructions: "Test instructions",
			ownerID:      1,
			mockSetup: func() {
				mockOpenAIClient.EXPECT().CreateVectorStore("1 - TestBot").Return("vsID", nil)
				mockOpenAIClient.EXPECT().CreateAssistant("1 - TestBot", "Test instructions", "vsID", nil).Return("", errors.New("assistant error"))
			},
			expectError: true,
			errorMsg:    "assistant error",
			expectedID:  0,
		},
		{
			name:         "Repository save error",
			nameArg:      "TestBot",
			description:  "A test bot",
			instructions: "Test instructions",
			ownerID:      1,
			mockSetup: func() {
				mockOpenAIClient.EXPECT().CreateVectorStore("1 - TestBot").Return("vsID", nil)
				mockOpenAIClient.EXPECT().CreateAssistant("1 - TestBot", "Test instructions", "vsID", nil).Return("assistantID", nil)
				mockRepo.EXPECT().SaveChatBot(ctx, gomock.Any()).Return(int64(0), errors.New("repository error"))
			},
			expectError: true,
			errorMsg:    "repository error",
			expectedID:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			id, err := u.SaveChatBot(ctx, tt.nameArg, tt.description, tt.instructions, tt.ownerID)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, int64(0), id)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

func TestUsecase_UpdateChatBot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBotRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:               logger,
		ChatBotRepository: mockRepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name         string
		id           int64
		ownerID      int64
		nameArg      string
		description  string
		instructions string
		mockSetup    func()
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "Valid chat bot update",
			id:           1,
			ownerID:      1,
			nameArg:      "UpdatedBot",
			description:  "An updated bot",
			instructions: "Updated instructions",
			mockSetup: func() {
				mockRepo.EXPECT().UpdateChatBot(ctx, int64(1), int64(1), domain.ChatBot{
					Name:         "UpdatedBot",
					Description:  "An updated bot",
					Instructions: "Updated instructions",
				}).Return(nil)
			},
			expectError: false,
		},
		{
			name:         "Validation error",
			id:           0,
			ownerID:      1,
			nameArg:      "",
			description:  "An updated bot",
			instructions: "Updated instructions",
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectError: true,
			errorMsg:    "invalid chat-bot id provided 0",
		},
		{
			name:         "Repository update error",
			id:           1,
			ownerID:      1,
			nameArg:      "UpdatedBot",
			description:  "An updated bot",
			instructions: "Updated instructions",
			mockSetup: func() {
				mockRepo.EXPECT().UpdateChatBot(ctx, int64(1), int64(1), domain.ChatBot{
					Name:         "UpdatedBot",
					Description:  "An updated bot",
					Instructions: "Updated instructions",
				}).Return(errors.New("repository error"))
			},
			expectError: true,
			errorMsg:    "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := u.UpdateChatBot(ctx, tt.id, tt.ownerID, tt.nameArg, tt.description, tt.instructions)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUsecase_RemoveChatBot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBotRepository(ctrl)
	mockOpenAIClient := mocks.NewMockOpenAIClient(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:               logger,
		ChatBotRepository: mockRepo,
		openAIClient:      mockOpenAIClient,
	}

	ctx := context.TODO()

	tests := []struct {
		name        string
		id          int64
		ownerID     int64
		mockSetup   func()
		expectError bool
		errorMsg    string
	}{
		{
			name:    "Valid chat bot removal",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					AssistantID:   "assistantID",
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().DeleteAssistant("assistantID").Return(nil)
				mockOpenAIClient.EXPECT().DeleteVectorStore("vectorStoreID").Return(nil)
				mockRepo.EXPECT().RemoveChatBot(ctx, int64(1), int64(1)).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "Validation error",
			id:      0,
			ownerID: 1,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectError: true,
			errorMsg:    "invalid chat-bot id provided 0",
		},
		{
			name:    "ChatBot retrieval error",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{}, errors.New("repository error"))
			},
			expectError: true,
			errorMsg:    "repository error",
		},
		{
			name:    "Delete assistant error",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					AssistantID:   "assistantID",
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().DeleteAssistant("assistantID").Return(errors.New("delete assistant error"))
			},
			expectError: true,
			errorMsg:    "delete assistant error",
		},
		{
			name:    "Delete vector store error",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					AssistantID:   "assistantID",
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().DeleteAssistant("assistantID").Return(nil)
				mockOpenAIClient.EXPECT().DeleteVectorStore("vectorStoreID").Return(errors.New("delete vector store error"))
			},
			expectError: true,
			errorMsg:    "delete vector store error",
		},
		{
			name:    "Remove chat bot error",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					AssistantID:   "assistantID",
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().DeleteAssistant("assistantID").Return(nil)
				mockOpenAIClient.EXPECT().DeleteVectorStore("vectorStoreID").Return(nil)
				mockRepo.EXPECT().RemoveChatBot(ctx, int64(1), int64(1)).Return(errors.New("remove chat bot error"))
			},
			expectError: true,
			errorMsg:    "remove chat bot error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := u.RemoveChatBot(ctx, tt.id, tt.ownerID)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUsecase_StartChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBotRepository(ctrl)
	mockOpenAIClient := mocks.NewMockOpenAIClient(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:               logger,
		ChatBotRepository: mockRepo,
		openAIClient:      mockOpenAIClient,
	}

	ctx := context.TODO()

	tests := []struct {
		name        string
		chatBotID   int64
		ownerID     int64
		mockSetup   func()
		expectError bool
		errorMsg    string
		expectedID  string
	}{
		{
			name:      "Valid start chat",
			chatBotID: 1,
			ownerID:   1,
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{}, nil)
				mockOpenAIClient.EXPECT().CreateThread().Return("threadID", nil)
			},
			expectError: false,
			expectedID:  "threadID",
		},
		{
			name:      "Invalid chat bot ID",
			chatBotID: 0,
			ownerID:   1,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectError: true,
			errorMsg:    "invalid chat-bot id provided 0",
			expectedID:  "",
		},
		{
			name:      "Invalid owner ID",
			chatBotID: 1,
			ownerID:   0,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectError: true,
			errorMsg:    "invalid owner id provided 0",
			expectedID:  "",
		},
		{
			name:      "ChatBot retrieval error",
			chatBotID: 1,
			ownerID:   1,
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{}, errors.New("repository error"))
			},
			expectError: true,
			errorMsg:    "repository error",
			expectedID:  "",
		},
		{
			name:      "Create thread error",
			chatBotID: 1,
			ownerID:   1,
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{}, nil)
				mockOpenAIClient.EXPECT().CreateThread().Return("", errors.New("create thread error"))
			},
			expectError: true,
			errorMsg:    "create thread error",
			expectedID:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			threadID, err := u.StartChat(ctx, tt.chatBotID, tt.ownerID)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, "", threadID)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, threadID)
			}
		})
	}
}

func TestUsecase_SendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBotRepository(ctrl)
	mockOpenAIClient := mocks.NewMockOpenAIClient(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:               logger,
		ChatBotRepository: mockRepo,
		openAIClient:      mockOpenAIClient,
		UserRepository:    mockUserRepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name          string
		botID         int64
		ownerID       int64
		threadID      string
		userPrompt    string
		mockSetup     func(responseChan chan AssistantResponse)
		expectedError error
		expectedResp  string
	}{
		{
			name:       "Valid message send",
			botID:      1,
			ownerID:    1,
			threadID:   "thread1",
			userPrompt: "Hello",
			mockSetup: func(responseChan chan AssistantResponse) {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					AssistantID: "assistant1",
				}, nil)
				mockOpenAIClient.EXPECT().AddMessageToThread("thread1", "Hello").Return(nil)
				mockOpenAIClient.EXPECT().CreateRun("thread1", "assistant1").Return("run1", nil)
				mockOpenAIClient.EXPECT().GetRun("thread1", "run1").Return(runs.GetRunResponse{
					Status: runs.StatusCompleted,
				}, nil).AnyTimes()
				mockOpenAIClient.EXPECT().LatestAssistantResponse("thread1").Return("Hi there!", nil)
				mockUserRepo.EXPECT().UpdateMessagesLeft(ctx, int64(1)).Return(nil)
			},
			expectedError: nil,
			expectedResp:  "Hi there!",
		},
		{
			name:       "Validation error",
			botID:      0,
			ownerID:    1,
			threadID:   "thread1",
			userPrompt: "Hello",
			mockSetup: func(responseChan chan AssistantResponse) {
				// No need to setup the mock since the validation will fail first
			},
			expectedError: errors.New("invalid chat-bot id provided 0"),
			expectedResp:  "",
		},
		{
			name:       "ChatBot retrieval error",
			botID:      1,
			ownerID:    1,
			threadID:   "thread1",
			userPrompt: "Hello",
			mockSetup: func(responseChan chan AssistantResponse) {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{}, errors.New("repository error"))
			},
			expectedError: errors.New("repository error"),
			expectedResp:  "",
		},
		{
			name:       "Add message to thread error",
			botID:      1,
			ownerID:    1,
			threadID:   "thread1",
			userPrompt: "Hello",
			mockSetup: func(responseChan chan AssistantResponse) {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					AssistantID: "assistant1",
				}, nil)
				mockOpenAIClient.EXPECT().AddMessageToThread("thread1", "Hello").Return(errors.New("add message error"))
			},
			expectedError: errors.New("add message error"),
			expectedResp:  "",
		},
		{
			name:       "Create run error",
			botID:      1,
			ownerID:    1,
			threadID:   "thread1",
			userPrompt: "Hello",
			mockSetup: func(responseChan chan AssistantResponse) {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					AssistantID: "assistant1",
				}, nil)
				mockOpenAIClient.EXPECT().AddMessageToThread("thread1", "Hello").Return(nil)
				mockOpenAIClient.EXPECT().CreateRun("thread1", "assistant1").Return("", errors.New("create run error"))
			},
			expectedError: errors.New("create run error"),
			expectedResp:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseChan := make(chan AssistantResponse)
			tt.mockSetup(responseChan)

			go u.SendMessage(ctx, tt.botID, tt.ownerID, tt.threadID, tt.userPrompt, responseChan)

			select {
			case response := <-responseChan:
				if tt.expectedError != nil {
					require.Error(t, response.Err)
					assert.Contains(t, response.Err.Error(), tt.expectedError.Error())
				} else {
					require.NoError(t, response.Err)
					assert.Equal(t, tt.expectedResp, response.Content)
				}
			case <-time.After(600 * time.Second):
				t.Fatalf("test timed out")
			}
		})
	}
}
