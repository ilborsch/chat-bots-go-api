package usecase

import (
	"chat-bots-api/domain"
	"chat-bots-api/internal/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
)

func TestUsecase_File(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mocks.NewMockFileRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:            logger,
		FileRepository: mockFileRepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name         string
		id           int64
		ownerID      int64
		mockSetup    func()
		expectedFile domain.File
		expectError  bool
		errorMsg     string
	}{
		{
			name:    "Valid file retrieval",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockFileRepo.EXPECT().File(ctx, int64(1), int64(1)).Return(domain.File{ID: 1, OwnerID: 1}, nil)
			},
			expectedFile: domain.File{ID: 1, OwnerID: 1},
			expectError:  false,
		},
		{
			name:    "Validation error",
			id:      0,
			ownerID: 1,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectedFile: domain.File{},
			expectError:  true,
			errorMsg:     "invalid file id provided 0",
		},
		{
			name:    "File retrieval error",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockFileRepo.EXPECT().File(ctx, int64(1), int64(1)).Return(domain.File{}, errors.New("repository error"))
			},
			expectedFile: domain.File{},
			expectError:  true,
			errorMsg:     "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			file, err := u.File(ctx, tt.id, tt.ownerID)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, tt.expectedFile, file)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedFile, file)
			}
		})
	}
}

func TestUsecase_ChatBotFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mocks.NewMockFileRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:            logger,
		FileRepository: mockFileRepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name          string
		chatBotID     int64
		ownerID       int64
		mockSetup     func()
		expectedFiles []domain.File
		expectError   bool
		errorMsg      string
	}{
		{
			name:      "Valid chat bot files retrieval",
			chatBotID: 1,
			ownerID:   1,
			mockSetup: func() {
				mockFileRepo.EXPECT().ChatBotFiles(ctx, int64(1), int64(1)).Return([]domain.File{{ID: 1, ChatBotID: 1}}, nil)
			},
			expectedFiles: []domain.File{{ID: 1, ChatBotID: 1}},
			expectError:   false,
		},
		{
			name:      "Validation error",
			chatBotID: 0,
			ownerID:   1,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectedFiles: nil,
			expectError:   true,
			errorMsg:      "invalid chat-bot id provided 0",
		},
		{
			name:      "File retrieval error",
			chatBotID: 1,
			ownerID:   1,
			mockSetup: func() {
				mockFileRepo.EXPECT().ChatBotFiles(ctx, int64(1), int64(1)).Return(nil, errors.New("repository error"))
			},
			expectedFiles: nil,
			expectError:   true,
			errorMsg:      "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			files, err := u.ChatBotFiles(ctx, tt.chatBotID, tt.ownerID)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, tt.expectedFiles, files)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedFiles, files)
			}
		})
	}
}

func TestUsecase_SaveFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBotRepository(ctrl)
	mockFileRepo := mocks.NewMockFileRepository(ctrl)
	mockOpenAIClient := mocks.NewMockOpenAIClient(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:               logger,
		ChatBotRepository: mockRepo,
		FileRepository:    mockFileRepo,
		openAIClient:      mockOpenAIClient,
	}

	ctx := context.TODO()

	tests := []struct {
		name        string
		chatBotID   int64
		ownerID     int64
		filename    string
		fileData    []byte
		mockSetup   func()
		expectedID  int64
		expectError bool
		errorMsg    string
	}{
		{
			name:      "Valid file save",
			chatBotID: 1,
			ownerID:   1,
			filename:  "test.txt",
			fileData:  []byte("test data"),
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().UploadFile("test.txt", []byte("test data")).Return("fileID", nil)
				mockOpenAIClient.EXPECT().AddVectorStoreFile("vectorStoreID", "fileID").Return(nil)
				mockFileRepo.EXPECT().SaveFile(ctx, gomock.Any(), int64(1)).Return(int64(1), nil)
			},
			expectedID:  1,
			expectError: false,
		},
		{
			name:      "Validation error for SaveFile",
			chatBotID: 1,
			ownerID:   1,
			filename:  "",
			fileData:  []byte("test data"),
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectedID:  0,
			expectError: true,
			errorMsg:    "empty filename is provided",
		},
		{
			name:      "Validation error for ChatBotID",
			chatBotID: 0,
			ownerID:   1,
			filename:  "test.txt",
			fileData:  []byte("test data"),
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectedID:  0,
			expectError: true,
			errorMsg:    "invalid chat-bot id provided 0",
		},
		{
			name:      "ChatBot retrieval error",
			chatBotID: 1,
			ownerID:   1,
			filename:  "test.txt",
			fileData:  []byte("test data"),
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{}, errors.New("repository error"))
			},
			expectedID:  0,
			expectError: true,
			errorMsg:    "repository error",
		},
		{
			name:      "File upload error",
			chatBotID: 1,
			ownerID:   1,
			filename:  "test.txt",
			fileData:  []byte("test data"),
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().UploadFile("test.txt", []byte("test data")).Return("", errors.New("upload error"))
			},
			expectedID:  0,
			expectError: true,
			errorMsg:    "upload error",
		},
		{
			name:      "Add file to vector store error",
			chatBotID: 1,
			ownerID:   1,
			filename:  "test.txt",
			fileData:  []byte("test data"),
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().UploadFile("test.txt", []byte("test data")).Return("fileID", nil)
				mockOpenAIClient.EXPECT().AddVectorStoreFile("vectorStoreID", "fileID").Return(errors.New("add vector store file error"))
			},
			expectedID:  0,
			expectError: true,
			errorMsg:    "add vector store file error",
		},
		{
			name:      "Save file to DB error",
			chatBotID: 1,
			ownerID:   1,
			filename:  "test.txt",
			fileData:  []byte("test data"),
			mockSetup: func() {
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().UploadFile("test.txt", []byte("test data")).Return("fileID", nil)
				mockOpenAIClient.EXPECT().AddVectorStoreFile("vectorStoreID", "fileID").Return(nil)
				mockFileRepo.EXPECT().SaveFile(ctx, gomock.Any(), int64(1)).Return(int64(0), errors.New("save file to DB error"))
			},
			expectedID:  0,
			expectError: true,
			errorMsg:    "save file to DB error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			id, err := u.SaveFile(ctx, tt.chatBotID, tt.ownerID, tt.filename, tt.fileData)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, tt.expectedID, id)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

func TestUsecase_RemoveFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBotRepository(ctrl)
	mockFileRepo := mocks.NewMockFileRepository(ctrl)
	mockOpenAIClient := mocks.NewMockOpenAIClient(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:               logger,
		ChatBotRepository: mockRepo,
		FileRepository:    mockFileRepo,
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
			name:    "Valid file removal",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockFileRepo.EXPECT().File(ctx, int64(1), int64(1)).Return(domain.File{
					ID:           1,
					ChatBotID:    1,
					OpenaiFileID: "fileID",
					FileSize:     100,
				}, nil)
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().DeleteVectorStoreFile("vectorStoreID", "fileID").Return(nil)
				mockOpenAIClient.EXPECT().DeleteFile("fileID").Return(nil)
				mockFileRepo.EXPECT().RemoveFile(ctx, int64(1), 100, int64(1)).Return(nil)
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
			errorMsg:    "invalid file id provided 0",
		},
		{
			name:    "File retrieval error",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockFileRepo.EXPECT().File(ctx, int64(1), int64(1)).Return(domain.File{}, errors.New("repository error"))
			},
			expectError: true,
			errorMsg:    "repository error",
		},
		{
			name:    "ChatBot retrieval error",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockFileRepo.EXPECT().File(ctx, int64(1), int64(1)).Return(domain.File{
					ID:           1,
					ChatBotID:    1,
					OpenaiFileID: "fileID",
					FileSize:     100,
				}, nil)
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{}, errors.New("repository error"))
			},
			expectError: true,
			errorMsg:    "repository error",
		},
		{
			name:    "Delete vector store file error",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockFileRepo.EXPECT().File(ctx, int64(1), int64(1)).Return(domain.File{
					ID:           1,
					ChatBotID:    1,
					OpenaiFileID: "fileID",
					FileSize:     100,
				}, nil)
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().DeleteVectorStoreFile("vectorStoreID", "fileID").Return(errors.New("delete vector store file error"))
			},
			expectError: true,
			errorMsg:    "delete vector store file error",
		},
		{
			name:    "Delete file error",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockFileRepo.EXPECT().File(ctx, int64(1), int64(1)).Return(domain.File{
					ID:           1,
					ChatBotID:    1,
					OpenaiFileID: "fileID",
					FileSize:     100,
				}, nil)
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().DeleteVectorStoreFile("vectorStoreID", "fileID").Return(nil)
				mockOpenAIClient.EXPECT().DeleteFile("fileID").Return(errors.New("delete file error"))
			},
			expectError: true,
			errorMsg:    "delete file error",
		},
		{
			name:    "Remove file from DB error",
			id:      1,
			ownerID: 1,
			mockSetup: func() {
				mockFileRepo.EXPECT().File(ctx, int64(1), int64(1)).Return(domain.File{
					ID:           1,
					ChatBotID:    1,
					OpenaiFileID: "fileID",
					FileSize:     100,
				}, nil)
				mockRepo.EXPECT().ChatBot(ctx, int64(1), int64(1)).Return(domain.ChatBot{
					VectorStoreID: "vectorStoreID",
				}, nil)
				mockOpenAIClient.EXPECT().DeleteVectorStoreFile("vectorStoreID", "fileID").Return(nil)
				mockOpenAIClient.EXPECT().DeleteFile("fileID").Return(nil)
				mockFileRepo.EXPECT().RemoveFile(ctx, int64(1), 100, int64(1)).Return(errors.New("remove file from DB error"))
			},
			expectError: true,
			errorMsg:    "remove file from DB error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := u.RemoveFile(ctx, tt.id, tt.ownerID)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
