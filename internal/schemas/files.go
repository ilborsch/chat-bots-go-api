package schemas

type ChatBotFilesRequest struct {
	ChatBotID int64 `json:"chat_bot_id"`
}

type File struct {
	ID        int64  `json:"id"`
	Filename  string `json:"filename"`
	ChatBotID int64  `json:"chat_bot_id"`
}

type ChatBotFiles struct {
	Files []File `json:"files"`
}

type SaveFileResponse struct {
	FileID int64 `json:"file_id"`
}

type RemoveFileRequest struct {
	FileID int64 `json:"file_id"`
}

type RemoveFileResponse struct {
	Success bool `json:"success"`
}
