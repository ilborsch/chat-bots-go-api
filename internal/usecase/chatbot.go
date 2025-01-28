package usecase

import (
	"chat-bots-api/domain"
	validate "chat-bots-api/internal/usecase/validators"
	"context"
	"errors"
	"fmt"
	"github.com/ilborsch/openai-go/openai/assistants/runs"
	"time"
)

var (
	AssistantResponseTimeoutError = errors.New("timeout error while waiting for assistant response")
)

type AssistantResponse struct {
	Err     error
	Content string
}

type ChatBotUsecase interface {
	ChatBot(ctx context.Context, id, ownerID int64) (domain.ChatBot, error)
	UserChatBots(ctx context.Context, ownerID int64) ([]domain.ChatBot, error)
	SaveChatBot(ctx context.Context, name, description, instructions string, ownerID int64) (id int64, err error)
	UpdateChatBot(ctx context.Context, id, ownerID int64, name, description, instructions string) error
	RemoveChatBot(ctx context.Context, id, ownerID int64) error
	StartChat(ctx context.Context, chatBotID, ownerID int64) (string, error)
	SendMessage(
		ctx context.Context,
		botID int64,
		ownerID int64,
		threadID string,
		userPrompt string,
		responseChan chan AssistantResponse,
	)
}

func (u *Usecase) ChatBot(ctx context.Context, id, ownerID int64) (domain.ChatBot, error) {
	if err := validate.ChatBotID(id); err != nil {
		u.log.Error("get chat-bot validation error: " + err.Error())
		return domain.ChatBot{}, err
	}
	if err := validate.OwnerID(ownerID); err != nil {
		u.log.Error("get chat-bot validation error: " + err.Error())
		return domain.ChatBot{}, err
	}
	chatBot, err := u.ChatBotRepository.ChatBot(ctx, id, ownerID)
	if err != nil {
		u.log.Error("error retrieving chat-bot: " + err.Error())
		return domain.ChatBot{}, err
	}
	return chatBot, nil
}

func (u *Usecase) UserChatBots(ctx context.Context, ownerID int64) ([]domain.ChatBot, error) {
	if err := validate.OwnerID(ownerID); err != nil {
		u.log.Error("get chat-bots validation error: " + err.Error())
		return nil, err
	}
	chatBots, err := u.ChatBotRepository.UserChatBots(ctx, ownerID)
	if err != nil {
		u.log.Error("error retrieving user chat-bots: " + err.Error())
		return nil, err
	}
	return chatBots, nil
}

func (u *Usecase) SaveChatBot(ctx context.Context, name, description, instructions string, ownerID int64) (id int64, err error) {
	chatBot := domain.ChatBot{
		Name:         name,
		Description:  description,
		Instructions: instructions,
		OwnerID:      ownerID,
	}
	if err := validate.SaveChatBot(chatBot); err != nil {
		u.log.Error("save chat-bot validation error: " + err.Error())
		return 0, err
	}

	client := u.openAIClient
	vsName := fmt.Sprintf("%v - %s", chatBot.OwnerID, chatBot.Name)
	vsID, err := client.CreateVectorStore(vsName)
	if err != nil {
		u.log.Error("error creating openai vector store: " + err.Error())
		return 0, err
	}

	assistantName := fmt.Sprintf("%v - %s", chatBot.OwnerID, chatBot.Name)
	assistantID, err := client.CreateAssistant(assistantName, chatBot.Instructions, vsID, nil)
	if err != nil {
		u.log.Error("error creating openai assistant: " + err.Error())
		return 0, err
	}

	chatBot.VectorStoreID = vsID
	chatBot.AssistantID = assistantID

	id, err = u.ChatBotRepository.SaveChatBot(ctx, chatBot)
	if err != nil {
		u.log.Error("error inserting chat-bot: " + err.Error())
		return 0, err
	}
	return id, nil
}

func (u *Usecase) UpdateChatBot(ctx context.Context, id, ownerID int64, name, description, instructions string) error {
	newChatBot := domain.ChatBot{
		Name:         name,
		Description:  description,
		Instructions: instructions,
	}
	if err := validate.UpdateChatBot(id, ownerID, newChatBot); err != nil {
		u.log.Error("update chat-bot validation error: " + err.Error())
		return err
	}

	if err := u.ChatBotRepository.UpdateChatBot(ctx, id, ownerID, newChatBot); err != nil {
		u.log.Error("error updating chat-bot: " + err.Error())
		return err
	}
	return nil
}

func (u *Usecase) RemoveChatBot(ctx context.Context, id, ownerID int64) error {
	if err := validate.ChatBotID(id); err != nil {
		u.log.Error("remove chat-bot validation error: " + err.Error())
		return err
	}

	chatBot, err := u.ChatBotRepository.ChatBot(ctx, id, ownerID)
	if err != nil {
		u.log.Error("error retrieving chat-bot: " + err.Error())
		return err
	}

	client := u.openAIClient
	assistantID := chatBot.AssistantID
	vectorStoreID := chatBot.VectorStoreID

	if err := client.DeleteAssistant(assistantID); err != nil {
		u.log.Error("error removing openai assistant: " + err.Error())
		return err
	}

	if err := client.DeleteVectorStore(vectorStoreID); err != nil {
		u.log.Error("error removing openai vector store: " + err.Error())
		return err
	}

	if err := u.ChatBotRepository.RemoveChatBot(ctx, id, ownerID); err != nil {
		u.log.Error("error removing chat-bot: " + err.Error())
		return err
	}
	return nil
}

func (u *Usecase) StartChat(ctx context.Context, chatBotID, ownerID int64) (string, error) {
	if err := validate.ChatBotID(chatBotID); err != nil {
		u.log.Error("get chat-bot validation error: " + err.Error())
		return "", err
	}
	if err := validate.OwnerID(ownerID); err != nil {
		u.log.Error("get chat-bot validation error: " + err.Error())
		return "", err
	}

	_, err := u.ChatBotRepository.ChatBot(ctx, chatBotID, ownerID)
	if err != nil {
		u.log.Error("error retrieving chat-bot: " + err.Error())
		return "", err
	}

	client := u.openAIClient
	threadID, err := client.CreateThread()
	if err != nil {
		u.log.Error("error creating thread: " + err.Error())
		return "", err
	}
	return threadID, nil
}

func (u *Usecase) SendMessage(
	ctx context.Context,
	botID int64,
	ownerID int64,
	threadID string,
	userPrompt string,
	responseChan chan AssistantResponse,
) {
	if err := validate.SendMessage(botID, threadID, userPrompt); err != nil {
		u.log.Error("validation error: " + err.Error())
		closeResponseChan(responseChan, err)
		return
	}

	chatBot, err := u.ChatBotRepository.ChatBot(ctx, botID, ownerID)
	if err != nil {
		u.log.Error("error retrieving chat-bot: " + err.Error())
		closeResponseChan(responseChan, err)
		return
	}

	assistantID := chatBot.AssistantID
	client := u.openAIClient

	if err := client.AddMessageToThread(threadID, userPrompt); err != nil {
		u.log.Error("error adding user message to a thread: " + err.Error())
		closeResponseChan(responseChan, err)
		return
	}

	runID, err := client.CreateRun(threadID, assistantID)
	if err != nil {
		u.log.Error("error adding user message to a thread: " + err.Error())
		closeResponseChan(responseChan, err)
		return
	}

	const poolInterval = 200 * time.Millisecond
	const maxPoolsAmount = 5000 / 200 // 5000 milliseconds timeout / 200 milliseconds per pool

	poolsAmount := 0
	runStatus := ""
	for poolsAmount < maxPoolsAmount {
		run, err := client.GetRun(threadID, runID)
		if err != nil {
			u.log.Error("error retrieving run status: " + err.Error())
			closeResponseChan(responseChan, err)
			return
		}

		runStatus = run.Status
		if runStatus == runs.StatusCompleted {
			assistantResponse, err := client.LatestAssistantResponse(threadID)
			if err != nil {
				u.log.Error("error retrieving assistant response: " + err.Error())
				closeResponseChan(responseChan, err)
				return
			}
			response := AssistantResponse{
				Content: assistantResponse,
				Err:     nil,
			}
			responseChan <- response
			close(responseChan)

			go func() {
				err := u.UserRepository.UpdateMessagesLeft(ctx, ownerID)
				if err != nil {
					u.log.Error("error updating user messages amount: " + err.Error())
					return
				}
			}()
			return
		}
		poolsAmount++
	}
	closeResponseChan(responseChan, AssistantResponseTimeoutError)
}

func closeResponseChan(responseChan chan AssistantResponse, err error) {
	responseChan <- AssistantResponse{
		Err:     err,
		Content: "",
	}
	close(responseChan)
}
