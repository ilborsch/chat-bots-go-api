package validate

import (
	"chat-bots-api/domain"
	"errors"
	"fmt"
)

func ChatBotID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid chat-bot id provided %v", id)
	}
	return nil
}

func OwnerID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid owner id provided %v", id)
	}
	return nil
}

func SaveChatBot(bot domain.ChatBot) error {
	if bot.Name == "" {
		return errors.New("empty chat-bot name provided")
	}
	return nil
}

func UpdateChatBot(id, ownerID int64, bot domain.ChatBot) error {
	if err := ChatBotID(id); err != nil {
		return err
	}
	if err := OwnerID(ownerID); err != nil {
		return err
	}
	if bot.Name == "" {
		return errors.New("cannot clear chat-bot name")
	}
	return nil
}

func SendMessage(botID int64, threadID, userPrompt string) error {
	if err := ChatBotID(botID); err != nil {
		return err
	}
	if threadID == "" {
		return errors.New("empty thread id passed")
	}
	if userPrompt == "" {
		return errors.New("cannot send blank user prompt")
	}
	return nil
}
