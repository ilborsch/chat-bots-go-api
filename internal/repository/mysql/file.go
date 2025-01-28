package mysql

import (
	"chat-bots-api/domain"
	"context"
	"errors"
	"fmt"
)

var (
	ErrFileAlreadyExists = errors.New("user already exists")
	ErrFileNotExists     = errors.New("user does not exist")
)

type FileRepository struct {
	*MySQL
}

func NewFileRepository(m *MySQL) *FileRepository {
	return &FileRepository{
		MySQL: m,
	}
}

func (m *FileRepository) File(ctx context.Context, id, ownerID int64) (domain.File, error) {
	const query = `
		SELECT * FROM file 
		WHERE id = ? AND owner_id = ?
	`

	res := m.db.QueryRowContext(ctx, query, id, ownerID)
	var file domain.File
	if err := res.Scan(
		&file.ID,
		&file.OwnerID,
		&file.ChatBotID,
		&file.OpenaiFileID,
		&file.Filename,
	); err != nil {
		return domain.File{}, err
	}
	return file, nil
}

func (m *FileRepository) ChatBotFiles(ctx context.Context, chatBotID, ownerID int64) ([]domain.File, error) {
	const query = `
		SELECT * FROM file 
		WHERE chat_bot_id = ? AND owner_id = ?
	`

	rows, err := m.db.QueryContext(ctx, query, chatBotID, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []domain.File
	for rows.Next() {
		var file domain.File
		if err := rows.Scan(
			&file.ID,
			&file.ChatBotID,
			&file.OwnerID,
			&file.OpenaiFileID,
			&file.Filename,
		); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (m *FileRepository) SaveFile(ctx context.Context, file domain.File, ownerID int64) (id int64, err error) {
	const saveFileQuery = `
		INSERT INTO file(chat_bot_id, owner_id, openai_file_id, filename)
		VALUES (?, ?, ?, ?);
	`
	const updateUserQuery = `
		UPDATE user
		SET bytes_data_left = bytes_data_left - ?
		WHERE id = ?
	`

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(
		ctx,
		saveFileQuery,
		file.ChatBotID,
		ownerID,
		file.OpenaiFileID,
		file.Filename,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	res, err = tx.ExecContext(
		ctx,
		updateUserQuery,
		file.FileSize,
		ownerID,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if n == 0 {
		tx.Rollback()
		return 0, fmt.Errorf("invalid owner ID passed %v", ownerID)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (m *FileRepository) RemoveFile(ctx context.Context, id int64, fileSize int, ownerID int64) error {
	const removeFileQuery = `
		DELETE FROM file WHERE id = ? AND owner_id = ?
	`
	const updateUserQuery = `
		UPDATE user
		SET bytes_data_left = bytes_data_left + ?
		WHERE id = ?
	`

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, removeFileQuery, id, ownerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, updateUserQuery, fileSize, ownerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
