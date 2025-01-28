package sso

import (
	"chat-bots-api/internal/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/ilborsch/sso-proto/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientRepository_IsAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthClient(ctrl)
	ctx := context.Background()
	request := &sso.IsAdminRequest{}
	expectedResponse := &sso.IsAdminResponse{IsAdmin: true}

	mockAuthClient.EXPECT().IsAdmin(ctx, request, gomock.Any()).Return(expectedResponse, nil)

	repo := &ClientRepository{client: mockAuthClient}
	isAdmin, err := repo.IsAdmin(ctx, request)

	require.NoError(t, err)
	assert.Equal(t, isAdmin, expectedResponse.IsAdmin)
}

func TestClientRepository_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthClient(ctrl)
	ctx := context.Background()
	request := &sso.LoginRequest{}
	expectedResponse := &sso.LoginResponse{Token: "test-token"}

	mockAuthClient.EXPECT().Login(ctx, request, gomock.Any()).Return(expectedResponse, nil)

	repo := &ClientRepository{client: mockAuthClient}
	token, err := repo.Login(ctx, request)

	require.NoError(t, err)
	assert.Equal(t, token, expectedResponse.Token)
}

func TestClientRepository_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthClient(ctrl)
	ctx := context.Background()
	request := &sso.RegisterRequest{}
	expectedResponse := &sso.RegisterResponse{UserId: 12345}

	mockAuthClient.EXPECT().Register(ctx, request, gomock.Any()).Return(expectedResponse, nil)

	repo := &ClientRepository{client: mockAuthClient}
	userID, err := repo.Register(ctx, request)

	require.NoError(t, err)
	assert.Equal(t, userID, expectedResponse.UserId)
}
