package domain

type File struct {
	ID           int64
	ChatBotID    int64
	OwnerID      int64
	OpenaiFileID string
	Filename     string
	FileSize     int
}
