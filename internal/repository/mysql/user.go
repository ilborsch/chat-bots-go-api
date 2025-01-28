package mysql

import (
	"chat-bots-api/domain"
	"context"
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotExists     = errors.New("user does not exist")
)

type UserRepository struct {
	*MySQL
}

func NewUserRepository(m *MySQL) *UserRepository {
	return &UserRepository{
		MySQL: m,
	}
}

func (m *UserRepository) User(ctx context.Context, id int64) (domain.User, error) {
	const query = `
		SELECT * FROM user 
		WHERE id = ?
	`
	res := m.db.QueryRowContext(ctx, query, id)

	var user domain.User
	if err := res.Scan(
		&user.ID,
		&user.Email,
		&user.Plan,
		&user.PlanBoughtDate,
		&user.MessagesLeft,
		&user.BytesDataLeft,
		&user.BotsLeft,
	); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (m *UserRepository) UserByEmail(ctx context.Context, email string) (domain.User, error) {
	const query = `
		SELECT * FROM user 
		WHERE email = ?
	`
	res := m.db.QueryRowContext(ctx, query, email)

	var user domain.User
	if err := res.Scan(
		&user.ID,
		&user.Email,
		&user.Plan,
		&user.PlanBoughtDate,
		&user.MessagesLeft,
		&user.BytesDataLeft,
		&user.BotsLeft,
	); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (m *UserRepository) SaveUser(ctx context.Context, user domain.User) (id int64, err error) {
	const query = `
		INSERT INTO user(email, plan, plan_bought_date, messages_left, bytes_data_left, bots_left)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	res, err := m.db.ExecContext(
		ctx,
		query,
		user.Email,
		user.Plan,
		user.PlanBoughtDate,
		user.MessagesLeft,
		user.BytesDataLeft,
		user.BotsLeft,
	)
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *UserRepository) UpdatePlan(ctx context.Context, id int64, newUser domain.User) error {
	const query = `
		UPDATE user
		SET plan = ?, plan_bought_date = ?, messages_left = ?, bytes_data_left = ?, bots_left = ?
		WHERE id = ?
	`
	res, err := m.db.ExecContext(
		ctx,
		query,
		newUser.Plan,
		newUser.PlanBoughtDate,
		newUser.MessagesLeft,
		newUser.BytesDataLeft,
		newUser.BotsLeft,
		id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotExists
	}
	return nil
}

func (m *UserRepository) UpdateMessagesLeft(ctx context.Context, id int64) error {
	const query = `
		UPDATE user
		SET messages_left = messages_left - 1
		WHERE id = ?
	`
	res, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotExists
	}
	return nil
}
