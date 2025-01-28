package validate

import (
	"chat-bots-api/domain"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChatBotID(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		expectError bool
	}{
		{"Valid ID", 1, false},
		{"Zero ID", 0, true},
		{"Negative ID", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ChatBotID(tt.id)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid chat-bot id provided")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOwnerID(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		expectError bool
	}{
		{"Valid ID", 1, false},
		{"Zero ID", 0, true},
		{"Negative ID", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OwnerID(tt.id)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid owner id provided")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSaveChatBot(t *testing.T) {
	tests := []struct {
		name        string
		bot         domain.ChatBot
		expectError bool
	}{
		{"Valid chat-bot", domain.ChatBot{Name: "Test Bot"}, false},
		{"Empty chat-bot name", domain.ChatBot{Name: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveChatBot(tt.bot)
			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, errors.New("empty chat-bot name provided"), err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateChatBot(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		ownerID     int64
		bot         domain.ChatBot
		expectError bool
		errorMsg    string
	}{
		{"Valid update", 1, 1, domain.ChatBot{Name: "Updated Bot"}, false, ""},
		{"Invalid chat-bot ID", 0, 1, domain.ChatBot{Name: "Updated Bot"}, true, "invalid chat-bot id provided"},
		{"Invalid owner ID", 1, 0, domain.ChatBot{Name: "Updated Bot"}, true, "invalid owner id provided"},
		{"Empty chat-bot name", 1, 1, domain.ChatBot{Name: ""}, true, "cannot clear chat-bot name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateChatBot(tt.id, tt.ownerID, tt.bot)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSendMessage(t *testing.T) {
	tests := []struct {
		name        string
		botID       int64
		threadID    string
		userPrompt  string
		expectError bool
		errorMsg    string
	}{
		{"Valid message", 1, "thread1", "Hello", false, ""},
		{"Invalid bot ID", 0, "thread1", "Hello", true, "invalid chat-bot id provided"},
		{"Empty thread ID", 1, "", "Hello", true, "empty thread id passed"},
		{"Empty user prompt", 1, "thread1", "", true, "cannot send blank user prompt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendMessage(tt.botID, tt.threadID, tt.userPrompt)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
