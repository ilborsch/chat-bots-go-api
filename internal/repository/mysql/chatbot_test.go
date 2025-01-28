package mysql

import (
	"chat-bots-api/domain"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestChatBotRepository(t *testing.T) (*ChatBotRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	m := &MySQL{db: db}
	repo := NewChatBotRepository(m)
	return repo, mock
}

func TestChatBotRepository_ChatBot(t *testing.T) {
	repo, mock := newTestChatBotRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	id, ownerID := int64(1), int64(1)

	const query = `SELECT \* FROM chat_bot WHERE id = \? AND owner_id = \?`
	rows := sqlmock.NewRows(
		[]string{"id", "assistant_id", "vector_store_id", "owner_id", "name", "description", "instructions"},
	).AddRow(1, "aID", "vsID", 1, "TestBot", "A test bot", "Instructions")

	mock.ExpectQuery(query).WithArgs(id, ownerID).WillReturnRows(rows)

	chatBot, err := repo.ChatBot(ctx, id, ownerID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), chatBot.ID)
	assert.Equal(t, "aID", chatBot.AssistantID)
	assert.Equal(t, "vsID", chatBot.VectorStoreID)
	assert.Equal(t, int64(1), chatBot.OwnerID)
	assert.Equal(t, "TestBot", chatBot.Name)
	assert.Equal(t, "A test bot", chatBot.Description)
	assert.Equal(t, "Instructions", chatBot.Instructions)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestChatBotRepository_UserChatBots(t *testing.T) {
	repo, mock := newTestChatBotRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	ownerID := int64(1)

	const query = `SELECT \* FROM chat_bot WHERE owner_id = \? ORDER BY name`
	rows := sqlmock.NewRows([]string{"id", "assistant_id", "vector_store_id", "owner_id", "name", "description", "instructions"}).
		AddRow(1, "aID", "vsID", 1, "TestBot1", "A test bot", "Instructions").
		AddRow(2, "aID", "vsID", 1, "TestBot2", "Another test bot", "More Instructions")

	mock.ExpectQuery(query).WithArgs(ownerID).WillReturnRows(rows)

	chatBots, err := repo.UserChatBots(ctx, ownerID)
	require.NoError(t, err)
	require.Len(t, chatBots, 2)

	assert.Equal(t, int64(1), chatBots[0].ID)
	assert.Equal(t, "TestBot1", chatBots[0].Name)
	assert.Equal(t, int64(2), chatBots[1].ID)
	assert.Equal(t, "TestBot2", chatBots[1].Name)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestChatBotRepository_SaveChatBot(t *testing.T) {
	repo, mock := newTestChatBotRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	chatBot := domain.ChatBot{
		AssistantID:   "aID",
		VectorStoreID: "vsID",
		OwnerID:       3,
		Name:          "TestBot",
		Description:   "A test bot",
		Instructions:  "Instructions",
	}

	const saveChatBotQuery = `INSERT INTO chat_bot\(assistant_id, vector_store_id, owner_id, name, description, instructions\) VALUES \(\?, \?, \?, \?, \?, \?\);`
	const updateUserQuery = `UPDATE users SET bots_left = bots_left - 1 WHERE id = \?`

	mock.ExpectBegin()

	mock.ExpectExec(saveChatBotQuery).
		WithArgs(chatBot.AssistantID, chatBot.VectorStoreID, chatBot.OwnerID, chatBot.Name, chatBot.Description, chatBot.Instructions).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(updateUserQuery).WithArgs(chatBot.OwnerID).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	id, err := repo.SaveChatBot(ctx, chatBot)
	require.NoError(t, err)
	assert.Equal(t, int64(1), id)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestChatBotRepository_UpdateChatBot(t *testing.T) {
	repo, mock := newTestChatBotRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	id, ownerID := int64(1), int64(1)
	chatBot := domain.ChatBot{
		Name:         "UpdatedBot",
		Description:  "Updated Description",
		Instructions: "Updated Instructions",
	}

	const query = `UPDATE chat_bot SET name = \?, description = \?, instructions = \? WHERE id = \? AND owner_id = \?`

	mock.ExpectExec(query).
		WithArgs(chatBot.Name, chatBot.Description, chatBot.Instructions, id, ownerID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateChatBot(ctx, id, ownerID, chatBot)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
