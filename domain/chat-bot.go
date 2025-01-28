package domain

type ChatBot struct {
	ID            int64
	AssistantID   string
	VectorStoreID string
	OwnerID       int64
	Name          string
	Description   string
	Instructions  string
}
