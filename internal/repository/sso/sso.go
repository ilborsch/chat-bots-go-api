package sso

import (
	"context"
	"github.com/ilborsch/sso-proto/gen/go/sso"
)

func (c *ClientRepository) IsAdmin(ctx context.Context, request *sso.IsAdminRequest) (bool, error) {
	response, err := c.client.IsAdmin(ctx, request)
	if err != nil {
		return false, err
	}
	return response.GetIsAdmin(), nil
}

func (c *ClientRepository) Login(ctx context.Context, request *sso.LoginRequest) (string, error) {
	response, err := c.client.Login(ctx, request)
	if err != nil {
		return "", err
	}
	return response.GetToken(), nil
}

func (c *ClientRepository) Register(ctx context.Context, request *sso.RegisterRequest) (int64, error) {
	response, err := c.client.Register(ctx, request)
	if err != nil {
		return 0, err
	}
	return response.GetUserId(), nil
}
