package usecase

import (
	validate "chat-bots-api/internal/usecase/validators"
	"context"
	"github.com/ilborsch/sso-proto/gen/go/sso"
)

type SSOUsecase interface {
	IsAdmin(ctx context.Context, uid int64) (bool, error)
	Login(ctx context.Context, email, password string, appID int) (string, error)
	Register(ctx context.Context, email, password string, uid int64) (int64, error)
}

func (u *Usecase) IsAdmin(ctx context.Context, uid int64) (bool, error) {
	if err := validate.UserID(uid); err != nil {
		u.log.Error("error validating uid: " + err.Error())
		return false, err
	}

	request := sso.IsAdminRequest{
		UidInApp: int32(uid),
	}

	isAdmin, err := u.SSORepository.IsAdmin(ctx, &request)
	if err != nil {
		u.log.Error("error retrieving response from grpc sso server: " + err.Error())
		return false, err
	}
	return isAdmin, nil
}

func (u *Usecase) Login(ctx context.Context, email, password string, appID int) (string, error) {
	if err := validate.Login(email, password); err != nil {
		u.log.Error("error validating login: " + err.Error())
		return "", err
	}

	request := sso.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    int32(appID),
	}

	token, err := u.SSORepository.Login(ctx, &request)
	if err != nil {
		u.log.Error("error retrieving response from grpc sso server: " + err.Error())
		return "", err
	}
	return token, nil
}

func (u *Usecase) Register(ctx context.Context, email, password string, uid int64) (int64, error) {
	if err := validate.Register(email, password, uid); err != nil {
		u.log.Error("error validating register: " + err.Error())
		return 0, err
	}

	request := sso.RegisterRequest{
		Email:    email,
		Password: password,
		UidInApp: int32(uid),
	}

	id, err := u.SSORepository.Register(ctx, &request)
	if err != nil {
		u.log.Error("error retrieving response from grpc sso server: " + err.Error())
		return 0, err
	}
	return id, nil
}
