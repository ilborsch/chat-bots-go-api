package mysql

import (
	"chat-bots-api/domain"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

// Create a new test repository
func newTestUserRepository(t *testing.T) (*UserRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	m := &MySQL{db: db}
	repo := NewUserRepository(m)
	return repo, mock
}

func TestUserRepository_User(t *testing.T) {
	repo, mock := newTestUserRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	id := int64(1)
	planBoughtDate := time.Date(2022, 10, 10, 0, 0, 0, 0, time.Local)

	const query = `SELECT \* FROM users WHERE id = \?`
	rows := sqlmock.NewRows([]string{"id", "email", "plan", "plan_bought_date", "messages_left", "bytes_data_left", "bots_left"}).
		AddRow(1, "test@example.com", domain.FreePlan, planBoughtDate, domain.FreeMessagesAmount, domain.FreeBytesDataAmount, domain.FreeBotsAmount)

	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(rows)

	user, err := repo.User(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, planBoughtDate, user.PlanBoughtDate)
	assert.Equal(t, domain.FreePlan, user.Plan)
	assert.Equal(t, domain.FreeMessagesAmount, user.MessagesLeft)
	assert.Equal(t, domain.FreeBytesDataAmount, user.BytesDataLeft)
	assert.Equal(t, domain.FreeBotsAmount, user.BotsLeft)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

// Test for the UserByEmail method
func TestUserRepository_UserByEmail(t *testing.T) {
	repo, mock := newTestUserRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	email := "test@example.com"
	planBoughtDate := time.Date(2022, 10, 10, 0, 0, 0, 0, time.Local)

	const query = `SELECT \* FROM users WHERE email = \?`
	rows := sqlmock.NewRows([]string{"id", "email", "plan", "plan_bought_date", "messages_left", "bytes_data_left", "bots_left"}).
		AddRow(1, "test@example.com", domain.FreePlan, planBoughtDate, domain.FreeMessagesAmount, domain.FreeBytesDataAmount, domain.FreeBotsAmount)

	mock.ExpectQuery(query).WithArgs(email).WillReturnRows(rows)

	user, err := repo.UserByEmail(ctx, email)
	require.NoError(t, err)
	require.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, planBoughtDate, user.PlanBoughtDate)
	assert.Equal(t, domain.FreePlan, user.Plan)
	assert.Equal(t, domain.FreeMessagesAmount, user.MessagesLeft)
	assert.Equal(t, domain.FreeBytesDataAmount, user.BytesDataLeft)
	assert.Equal(t, domain.FreeBotsAmount, user.BotsLeft)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestUserRepository_SaveUser(t *testing.T) {
	repo, mock := newTestUserRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	user := domain.User{
		Email:          "test@email.com",
		Plan:           domain.FreePlan,
		PlanBoughtDate: time.Now(),
		MessagesLeft:   domain.FreeMessagesAmount,
		BytesDataLeft:  domain.FreeBytesDataAmount,
		BotsLeft:       domain.FreeBotsAmount,
	}

	query := regexp.QuoteMeta(`
		INSERT INTO users(email, plan, plan_bought_date, messages_left, bytes_data_left, bots_left)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	mock.ExpectExec(query).
		WithArgs(user.Email, user.Plan, user.PlanBoughtDate, user.MessagesLeft, user.BytesDataLeft, user.BotsLeft).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := repo.SaveUser(ctx, user)
	if err != nil {
		t.Fatalf("error executing SaveUser: %v", err)
	}
	require.NoError(t, err)
	assert.Equal(t, int64(1), id)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestUserRepository_UpdatePlan(t *testing.T) {
	repo, mock := newTestUserRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	id := int64(1)
	newUser := domain.User{
		Plan:           domain.FreePlan,
		PlanBoughtDate: time.Now(),
		MessagesLeft:   domain.FreeMessagesAmount,
		BytesDataLeft:  domain.FreeBytesDataAmount,
		BotsLeft:       domain.FreeBotsAmount,
	}

	const query = `UPDATE users SET plan = \?, plan_bought_date = \?, messages_left = \?, bytes_data_left = \?, bots_left = \? WHERE id = \?`

	mock.ExpectExec(query).
		WithArgs(newUser.Plan, newUser.PlanBoughtDate, newUser.MessagesLeft, newUser.BytesDataLeft, newUser.BotsLeft, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdatePlan(ctx, id, newUser)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestUserRepository_UpdateMessagesLeft(t *testing.T) {
	repo, mock := newTestUserRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	id := int64(1)

	const query = `UPDATE users SET messages_left = messages_left - 1 WHERE id = \?`

	mock.ExpectExec(query).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateMessagesLeft(ctx, id)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
