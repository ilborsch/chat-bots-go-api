package validate

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		expectError bool
		errorMsg    string
	}{
		{"Valid login", "user@example.com", "password123", false, ""},
		{"Invalid email", "invalid-email", "password123", true, "invalid email format"},
		{"Short password", "user@example.com", "short", true, "password is too short"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Login(tt.email, tt.password)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		uid         int64
		expectError bool
		errorMsg    string
	}{
		{"Valid register", "user@example.com", "password123", 1, false, ""},
		{"Invalid user ID", "user@example.com", "password123", 0, true, "invalid owner id provided"},
		{"Invalid email", "invalid-email", "password123", 1, true, "invalid email format"},
		{"Short password", "user@example.com", "short", 1, true, "password is too short"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Register(tt.email, tt.password, tt.uid)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
