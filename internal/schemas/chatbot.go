package schemas

type ChatBot struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Instructions string `json:"instructions"`
}

type SaveChatBotRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Instructions string `json:"instructions"`
}

type SaveChatBotResponse struct {
	ChatBotID int64 `json:"chat_bot_id"`
}

type UpdateChatBotRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Instructions string `json:"instructions"`
}

type UpdateChatBotResponse struct {
	Success bool `json:"success"`
}

type RemoveChatBotResponse struct {
	Success bool `json:"success"`
}

type SendMessageResponse struct {
	Response string `json:"response"`
}

type SendMessageRequest struct {
	Message string `json:"message"`
}

type SendMessageError struct {
	Error string `json:"error"`
}

type UserChatBots struct {
	ChatBots []ChatBot `json:"chat_bots"`
}
