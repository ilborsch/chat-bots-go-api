package usecase

import (
	"chat-bots-api/domain"
	validate "chat-bots-api/internal/usecase/validators"
	"context"
	"fmt"
)

type FileUsecase interface {
	File(ctx context.Context, id, ownerID int64) (domain.File, error)
	ChatBotFiles(ctx context.Context, chatBotID, ownerID int64) ([]domain.File, error)
	SaveFile(ctx context.Context, chatBotID int64, ownerID int64, filename string, fileData []byte) (id int64, err error)
	RemoveFile(ctx context.Context, id int64, ownerID int64) error
}

func (u *Usecase) File(ctx context.Context, id, ownerID int64) (domain.File, error) {
	if err := validate.File(id, ownerID); err != nil {
		u.log.Error("error validating file: " + err.Error())
		return domain.File{}, err
	}

	file, err := u.FileRepository.File(ctx, id, ownerID)
	if err != nil {
		u.log.Error("error retrieving file: " + err.Error())
		return domain.File{}, err
	}
	return file, nil
}

func (u *Usecase) ChatBotFiles(ctx context.Context, chatBotID, ownerID int64) ([]domain.File, error) {
	if err := validate.ChatBotID(chatBotID); err != nil {
		u.log.Error("error validating chat-bot id: " + err.Error())
		return nil, err
	}
	if err := validate.OwnerID(ownerID); err != nil {
		u.log.Error("error validating owner id: " + err.Error())
		return nil, err
	}

	if _, err := u.ChatBotRepository.ChatBot(ctx, chatBotID, ownerID); err != nil {
		u.log.Error(fmt.Sprintf("no chat-bot exists with id %v", chatBotID))
		return nil, err
	}

	files, err := u.FileRepository.ChatBotFiles(ctx, chatBotID, ownerID)
	if err != nil {
		u.log.Error("error retrieving files: " + err.Error())
		return nil, err
	}
	return files, nil
}

func (u *Usecase) SaveFile(
	ctx context.Context,
	chatBotID int64,
	ownerID int64,
	filename string,
	fileData []byte,
) (id int64, err error) {
	if err := validate.SaveFile(filename, fileData); err != nil {
		u.log.Error("error validating file: " + err.Error())
		return 0, err
	}
	if err := validate.ChatBotID(chatBotID); err != nil {
		u.log.Error("error validating chat-bot id: " + err.Error())
		return 0, err
	}
	if err := validate.OwnerID(ownerID); err != nil {
		u.log.Error("error validating owner id: " + err.Error())
		return 0, err
	}

	chatBot, err := u.ChatBotRepository.ChatBot(ctx, chatBotID, ownerID)
	if err != nil {
		u.log.Error("error retrieving chat-bot: " + err.Error())
		return 0, err
	}

	client := u.openAIClient
	fileID, err := client.UploadFile(filename, fileData)
	if err != nil {
		u.log.Error("error uploading file to openai: " + err.Error())
		return 0, err
	}

	vectorStoreID := chatBot.VectorStoreID
	if err := client.AddVectorStoreFile(vectorStoreID, fileID); err != nil {
		u.log.Error("error adding file to openai vector store: " + err.Error())
		return 0, err
	}

	file := domain.File{
		ChatBotID:    chatBotID,
		OpenaiFileID: fileID,
		Filename:     filename,
		FileSize:     len(fileData),
	}

	id, err = u.FileRepository.SaveFile(ctx, file, ownerID)
	if err != nil {
		u.log.Error("error creating file object: " + err.Error())
		return 0, err
	}
	return id, nil
}

func (u *Usecase) RemoveFile(ctx context.Context, id int64, ownerID int64) error {
	if err := validate.File(id, ownerID); err != nil {
		u.log.Error("error validating file id: " + err.Error())
		return err
	}

	file, err := u.FileRepository.File(ctx, id, ownerID)
	if err != nil {
		u.log.Error("error retrieving file: " + err.Error())
		return err
	}

	chatBot, err := u.ChatBotRepository.ChatBot(ctx, file.ChatBotID, ownerID)
	if err != nil {
		u.log.Error("error retrieving file's chat-bot: " + err.Error())
		return err
	}

	client := u.openAIClient
	vectorStoreID := chatBot.VectorStoreID

	if err := client.DeleteVectorStoreFile(vectorStoreID, file.OpenaiFileID); err != nil {
		u.log.Error("error deleting file from openai vector store: " + err.Error())
		return err
	}

	if err := client.DeleteFile(file.OpenaiFileID); err != nil {
		u.log.Error("error deleting file from openai: " + err.Error())
		return err
	}

	if err := u.FileRepository.RemoveFile(ctx, id, file.FileSize, ownerID); err != nil {
		u.log.Error("error removing file from db: " + err.Error())
		return err
	}
	return nil
}
