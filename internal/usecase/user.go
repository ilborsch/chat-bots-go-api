package usecase

import (
	"chat-bots-api/domain"
	validate "chat-bots-api/internal/usecase/validators"
	"context"
	"time"
)

type UserUsecase interface {
	User(ctx context.Context, id int64) (domain.User, error)
	UserByEmail(ctx context.Context, email string) (domain.User, error)
	SaveUser(ctx context.Context, email, password, plan string) (id int64, err error)
	UpdatePlan(ctx context.Context, id int64, plan string) error
}

func (u *Usecase) User(ctx context.Context, id int64) (domain.User, error) {
	if err := validate.UserID(id); err != nil {
		u.log.Error("error validating user id: " + err.Error())
		return domain.User{}, err
	}

	user, err := u.UserRepository.User(ctx, id)
	if err != nil {
		u.log.Error("error retrieving user from db: " + err.Error())
		return domain.User{}, err
	}
	return user, nil
}

func (u *Usecase) UserByEmail(ctx context.Context, email string) (domain.User, error) {
	if err := validate.UserEmail(email); err != nil {
		u.log.Error("error validating user email: " + err.Error())
		return domain.User{}, err
	}

	user, err := u.UserRepository.UserByEmail(ctx, email)
	if err != nil {
		u.log.Error("error retrieving user by email from db: " + err.Error())
		return domain.User{}, err
	}
	return user, nil
}

func (u *Usecase) SaveUser(ctx context.Context, email, password, plan string) (id int64, err error) {
	if err := validate.SaveUser(email, password, plan); err != nil {
		u.log.Error("error validating new user: " + err.Error())
		return 0, err
	}

	newUser := domain.User{
		Plan:  plan,
		Email: email,
	}

	setupUserPayloadByPlan(&newUser, newUser.Plan)
	newUser.PlanBoughtDate = time.Now()

	id, err = u.UserRepository.SaveUser(ctx, newUser)
	if err != nil {
		u.log.Error("error creating new user: " + err.Error())
		return 0, err
	}
	return id, err
}

func (u *Usecase) UpdatePlan(ctx context.Context, id int64, plan string) error {
	if err := validate.Plan(plan); err != nil {
		u.log.Error("error validating plan: " + err.Error())
		return err
	}
	newUser := domain.User{
		Plan: plan,
	}
	setupUserPayloadByPlan(&newUser, plan)
	newUser.PlanBoughtDate = time.Now()

	if err := u.UserRepository.UpdatePlan(ctx, id, newUser); err != nil {
		u.log.Error("error updating user plan: " + err.Error())
		return err
	}
	return nil
}

func setupUserPayloadByPlan(user *domain.User, plan string) {
	if plan == domain.FreePlan {
		user.BotsLeft = domain.FreeBotsAmount
		user.MessagesLeft = domain.FreeMessagesAmount
		user.BytesDataLeft = domain.FreeBytesDataAmount
		return
	}
	if plan == domain.BusinessPlan {
		user.BotsLeft = domain.BusinessBotsAmount
		user.MessagesLeft = domain.BusinessMessagesAmount
		user.BytesDataLeft = domain.BusinessBytesDataAmount
		return
	}
	if plan == domain.EnterprisePlan {
		user.BotsLeft = domain.EnterpriseBotsAmount
		user.MessagesLeft = domain.EnterpriseMessagesAmount
		user.BytesDataLeft = domain.EnterpriseBytesDataAmount
	}
}
