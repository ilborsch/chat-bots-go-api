package mysql

import (
	"chat-bots-api/domain"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func newTestFileRepository(t *testing.T) (*FileRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	m := &MySQL{db: db}
	repo := NewFileRepository(m)
	return repo, mock
}

func TestFileRepository_File(t *testing.T) {
	repo, mock := newTestFileRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	id, ownerID := int64(1), int64(1)

	const query = `SELECT \* FROM file WHERE id = \? AND owner_id = \?`
	rows := sqlmock.NewRows([]string{"id", "chat_bot_id", "openai_file_id", "filename"}).
		AddRow(1, 2, "file-id", "testfile.txt")

	mock.ExpectQuery(query).WithArgs(id, ownerID).WillReturnRows(rows)

	file, err := repo.File(ctx, id, ownerID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), file.ID)
	assert.Equal(t, int64(2), file.ChatBotID)
	assert.Equal(t, "file-id", file.OpenaiFileID)
	assert.Equal(t, "testfile.txt", file.Filename)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestFileRepository_ChatBotFiles(t *testing.T) {
	repo, mock := newTestFileRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	chatBotID, ownerID := int64(1), int64(1)

	const query = `SELECT \* FROM file WHERE chat_bot_id = \? AND owner_id = \?`
	rows := sqlmock.NewRows([]string{"id", "chat_bot_id", "openai_file_id", "filename"}).
		AddRow(1, chatBotID, "file-id-1", "testfile1.txt").
		AddRow(2, chatBotID, "file-id-2", "testfile2.txt")

	mock.ExpectQuery(query).WithArgs(chatBotID, ownerID).WillReturnRows(rows)

	files, err := repo.ChatBotFiles(ctx, chatBotID, ownerID)
	require.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Equal(t, int64(1), files[0].ID)
	assert.Equal(t, int64(2), files[1].ID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestFileRepository_SaveFile(t *testing.T) {
	repo, mock := newTestFileRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	file := domain.File{
		ChatBotID:    1,
		OpenaiFileID: "file-id",
		Filename:     "testfile.txt",
		FileSize:     1024,
	}
	ownerID := int64(1)

	const saveFileQuery = `INSERT INTO file\(chat_bot_id, openai_file_id, filename\) VALUES \(\?, \?, \?\)`
	const updateUserQuery = `UPDATE users SET bytes_data_left = bytes_data_left - \? WHERE id = \?`

	mock.ExpectBegin()

	mock.ExpectExec(saveFileQuery).
		WithArgs(file.ChatBotID, file.OpenaiFileID, file.Filename).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(updateUserQuery).
		WithArgs(file.FileSize, ownerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	id, err := repo.SaveFile(ctx, file, ownerID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), id)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestFileRepository_RemoveFile(t *testing.T) {
	repo, mock := newTestFileRepository(t)
	defer repo.MySQL.db.Close()

	ctx := context.Background()
	id, ownerID := int64(1), int64(1)
	fileSize := 1024

	const removeFileQuery = `DELETE FROM files WHERE id = \? AND owner_id = \?`
	const updateUserQuery = `UPDATE users SET bytes_data_left = bytes_data_left \+ \? WHERE id = \?`

	mock.ExpectBegin()

	mock.ExpectExec(removeFileQuery).
		WithArgs(id, ownerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(updateUserQuery).
		WithArgs(fileSize, ownerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repo.RemoveFile(ctx, id, fileSize, ownerID)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
