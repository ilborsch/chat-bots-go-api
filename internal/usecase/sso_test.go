package usecase

import (
	"chat-bots-api/internal/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/ilborsch/sso-proto/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
)

func TestUsecase_IsAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSSORepo := mocks.NewMockSSORepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:           logger,
		SSORepository: mockSSORepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name        string
		uid         int64
		mockSetup   func()
		expected    bool
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid is admin check",
			uid:  1,
			mockSetup: func() {
				mockSSORepo.EXPECT().IsAdmin(ctx, &sso.IsAdminRequest{UidInApp: int32(1)}).Return(true, nil)
			},
			expected:    true,
			expectError: false,
		},
		{
			name: "Validation error",
			uid:  0,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expected:    false,
			expectError: true,
			errorMsg:    "invalid owner id provided 0",
		},
		{
			name: "SSO repository error",
			uid:  1,
			mockSetup: func() {
				mockSSORepo.EXPECT().IsAdmin(ctx, &sso.IsAdminRequest{UidInApp: int32(1)}).Return(false, errors.New("repository error"))
			},
			expected:    false,
			expectError: true,
			errorMsg:    "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			isAdmin, err := u.IsAdmin(ctx, tt.uid)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, tt.expected, isAdmin)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, isAdmin)
			}
		})
	}
}

func TestUsecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSSORepo := mocks.NewMockSSORepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:           logger,
		SSORepository: mockSSORepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int
		mockSetup   func()
		expected    string
		expectError bool
		errorMsg    string
	}{
		{
			name:     "Valid login",
			email:    "test@example.com",
			password: "password",
			appID:    1,
			mockSetup: func() {
				mockSSORepo.EXPECT().Login(ctx, &sso.LoginRequest{
					Email:    "test@example.com",
					Password: "password",
					AppId:    int32(1),
				}).Return("token", nil)
			},
			expected:    "token",
			expectError: false,
		},
		{
			name:     "Validation error",
			email:    "",
			password: "password",
			appID:    1,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expected:    "",
			expectError: true,
			errorMsg:    "empty email provided",
		},
		{
			name:     "SSO repository error",
			email:    "test@example.com",
			password: "password",
			appID:    1,
			mockSetup: func() {
				mockSSORepo.EXPECT().Login(ctx, &sso.LoginRequest{
					Email:    "test@example.com",
					Password: "password",
					AppId:    int32(1),
				}).Return("", errors.New("repository error"))
			},
			expected:    "",
			expectError: true,
			errorMsg:    "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			token, err := u.Login(ctx, tt.email, tt.password, tt.appID)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, tt.expected, token)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, token)
			}
		})
	}
}

func TestUsecase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSSORepo := mocks.NewMockSSORepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	u := &Usecase{
		log:           logger,
		SSORepository: mockSSORepo,
	}

	ctx := context.TODO()

	tests := []struct {
		name        string
		email       string
		password    string
		uid         int64
		mockSetup   func()
		expected    int64
		expectError bool
		errorMsg    string
	}{
		{
			name:     "Valid register",
			email:    "test@example.com",
			password: "password",
			uid:      1,
			mockSetup: func() {
				mockSSORepo.EXPECT().Register(ctx, &sso.RegisterRequest{
					Email:    "test@example.com",
					Password: "password",
					UidInApp: int32(1),
				}).Return(int64(1), nil)
			},
			expected:    1,
			expectError: false,
		},
		{
			name:     "Validation error",
			email:    "",
			password: "password",
			uid:      1,
			mockSetup: func() {
				// No need to setup the mock since the validation will fail first
			},
			expected:    0,
			expectError: true,
			errorMsg:    "empty email provided",
		},
		{
			name:     "SSO repository error",
			email:    "test@example.com",
			password: "password",
			uid:      1,
			mockSetup: func() {
				mockSSORepo.EXPECT().Register(ctx, &sso.RegisterRequest{
					Email:    "test@example.com",
					Password: "password",
					UidInApp: int32(1),
				}).Return(int64(0), errors.New("repository error"))
			},
			expected:    0,
			expectError: true,
			errorMsg:    "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			id, err := u.Register(ctx, tt.email, tt.password, tt.uid)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, tt.expected, id)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, id)
			}
		})
	}
}
