package mysql

import (
	"chat-bots-api/domain"
	"context"
	"errors"
	"fmt"
)

var (
	ErrChatBotAlreadyExists = errors.New("chat bot already exists")
	ErrChatBotNotExists     = errors.New("chat bot does not exist")
)

type ChatBotRepository struct {
	*MySQL
}

func NewChatBotRepository(m *MySQL) *ChatBotRepository {
	return &ChatBotRepository{
		MySQL: m,
	}
}

func (m *ChatBotRepository) ChatBot(ctx context.Context, id, ownerID int64) (domain.ChatBot, error) {
	const query = `
		SELECT * FROM chat_bot 
		WHERE id = ? AND owner_id = ?
	`

	res := m.db.QueryRowContext(ctx, query, id, ownerID)
	var chatBot domain.ChatBot
	if err := res.Scan(
		&chatBot.ID,
		&chatBot.AssistantID,
		&chatBot.VectorStoreID,
		&chatBot.OwnerID,
		&chatBot.Name,
		&chatBot.Description,
		&chatBot.Instructions,
	); err != nil {
		return domain.ChatBot{}, err
	}
	return chatBot, nil
}

func (m *ChatBotRepository) UserChatBots(ctx context.Context, ownerID int64) ([]domain.ChatBot, error) {
	const query = `
		SELECT * FROM chat_bot
		WHERE owner_id = ?
		ORDER BY name
	`

	rows, err := m.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatBots []domain.ChatBot
	for rows.Next() {
		var chatBot domain.ChatBot
		if err := rows.Scan(
			&chatBot.ID,
			&chatBot.AssistantID,
			&chatBot.VectorStoreID,
			&chatBot.OwnerID,
			&chatBot.Name,
			&chatBot.Description,
			&chatBot.Instructions,
		); err != nil {
			return nil, err
		}
		chatBots = append(chatBots, chatBot)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return chatBots, nil
}

func (m *ChatBotRepository) SaveChatBot(ctx context.Context, chatBot domain.ChatBot) (id int64, err error) {
	const saveChatBotQuery = `
		INSERT INTO chat_bot(assistant_id, vector_store_id, owner_id, name, description, instructions)
		VALUES (?, ?, ?, ?, ?, ?);
	`
	const updateUserQuery = `
		UPDATE user
		SET bots_left = bots_left - 1
		WHERE id = ?
	`

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(
		ctx,
		saveChatBotQuery,
		chatBot.AssistantID,
		chatBot.VectorStoreID,
		chatBot.OwnerID,
		chatBot.Name,
		chatBot.Description,
		chatBot.Instructions,
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
		chatBot.OwnerID,
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
		return 0, fmt.Errorf("invalid owner id %v", chatBot.OwnerID)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ChatBotRepository) UpdateChatBot(ctx context.Context, id, ownerID int64, cb domain.ChatBot) error {
	const query = `
		UPDATE chat_bot
		SET name = ?, description = ?, instructions = ?
		WHERE id = ? AND owner_id = ?
	`
	res, err := m.db.ExecContext(ctx, query, cb.Name, cb.Description, cb.Instructions, id, ownerID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrChatBotNotExists
	}
	return nil
}

func (m *ChatBotRepository) RemoveChatBot(ctx context.Context, id int64, userID int64) error {
	const removeChatBotQuery = `
		DELETE FROM chat_bot WHERE id = ? AND owner_id = ?
	`
	const updateUserQuery = `
		UPDATE user
		SET bots_left = bots_left + 1
		WHERE id = ?
	`

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, removeChatBotQuery, id, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if n == 0 {
		tx.Rollback()
		return fmt.Errorf("no chat bot found with id %d and ownerID %d", id, userID)
	}

	res, err = tx.ExecContext(ctx, updateUserQuery, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	n, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if n == 0 {
		tx.Rollback()
		return fmt.Errorf("invalid user id %d", userID)
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
