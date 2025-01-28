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

func TestUsecase_User(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:            logger,
		UserRepository: mockUserRepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name        string
		id          int64
		mockSetup   func()
		expected    domain.User
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid user retrieval",
			id:   1,
			mockSetup: func() {
				mockUserRepo.EXPECT().User(ctx, int64(1)).Return(domain.User{ID: 1, Email: "test@example.com"}, nil)
			},
			expected:    domain.User{ID: 1, Email: "test@example.com"},
			expectError: false,
		},
		{
			name: "Validation error",
			id:   0,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expected:    domain.User{},
			expectError: true,
			errorMsg:    "invalid owner id provided 0",
		},
		{
			name: "User retrieval error",
			id:   1,
			mockSetup: func() {
				mockUserRepo.EXPECT().User(ctx, int64(1)).Return(domain.User{}, errors.New("repository error"))
			},
			expected:    domain.User{},
			expectError: true,
			errorMsg:    "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := u.User(ctx, tt.id)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, tt.expected, user)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, user)
			}
		})
	}
}

func TestUsecase_UserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:            logger,
		UserRepository: mockUserRepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name        string
		email       string
		mockSetup   func()
		expected    domain.User
		expectError bool
		errorMsg    string
	}{
		{
			name:  "Valid user retrieval by email",
			email: "test@example.com",
			mockSetup: func() {
				mockUserRepo.EXPECT().UserByEmail(ctx, "test@example.com").Return(domain.User{ID: 1, Email: "test@example.com"}, nil)
			},
			expected:    domain.User{ID: 1, Email: "test@example.com"},
			expectError: false,
		},
		{
			name:  "Validation error",
			email: "",
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expected:    domain.User{},
			expectError: true,
			errorMsg:    "empty email provided",
		},
		{
			name:  "User retrieval error",
			email: "test@example.com",
			mockSetup: func() {
				mockUserRepo.EXPECT().UserByEmail(ctx, "test@example.com").Return(domain.User{}, errors.New("repository error"))
			},
			expected:    domain.User{},
			expectError: true,
			errorMsg:    "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := u.UserByEmail(ctx, tt.email)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, tt.expected, user)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, user)
			}
		})
	}
}

func TestUsecase_SaveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:            logger,
		UserRepository: mockUserRepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name        string
		email       string
		password    string
		plan        string
		mockSetup   func()
		expectedID  int64
		expectError bool
		errorMsg    string
	}{
		{
			name:     "Valid user save",
			email:    "test@example.com",
			password: "password",
			plan:     domain.FreePlan,
			mockSetup: func() {
				mockUserRepo.EXPECT().SaveUser(ctx, gomock.Any()).Return(int64(1), nil)
			},
			expectedID:  1,
			expectError: false,
		},
		{
			name:     "Validation error",
			email:    "",
			password: "password",
			plan:     "Free",
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectedID:  0,
			expectError: true,
			errorMsg:    "invalid user plan provided",
		},
		{
			name:     "User save error",
			email:    "test@example.com",
			password: "password",
			plan:     domain.FreePlan,
			mockSetup: func() {
				mockUserRepo.EXPECT().SaveUser(ctx, gomock.Any()).Return(int64(0), errors.New("repository error"))
			},
			expectedID:  0,
			expectError: true,
			errorMsg:    "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			id, err := u.SaveUser(ctx, tt.email, tt.password, tt.plan)

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

func TestUsecase_UpdatePlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:            logger,
		UserRepository: mockUserRepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name        string
		id          int64
		plan        string
		mockSetup   func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid plan update",
			id:   1,
			plan: domain.FreePlan,
			mockSetup: func() {
				mockUserRepo.EXPECT().UpdatePlan(ctx, int64(1), gomock.Any()).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Validation error",
			id:   1,
			plan: "",
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expectError: true,
			errorMsg:    "invalid user plan provided",
		},
		{
			name: "Plan update error",
			id:   1,
			plan: domain.FreePlan,
			mockSetup: func() {
				mockUserRepo.EXPECT().UpdatePlan(ctx, int64(1), gomock.Any()).Return(errors.New("repository error"))
			},
			expectError: true,
			errorMsg:    "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := u.UpdatePlan(ctx, tt.id, tt.plan)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
