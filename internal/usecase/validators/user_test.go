package validate

import (
	"chat-bots-api/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserID(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		expectError bool
		errorMsg    string
	}{
		{"Valid user ID", 1, false, ""},
		{"Invalid user ID", 0, true, "invalid owner id provided 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UserID(tt.id)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError bool
		errorMsg    string
	}{
		{"Valid email", "user@example.com", false, ""},
		{"Empty email", "", true, "empty email provided"},
		{"Missing at symbol", "userexample.com", true, "invalid email format"},
		{"Missing dot", "user@examplecom", true, "invalid email format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UserEmail(tt.email)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPlan(t *testing.T) {
	tests := []struct {
		name        string
		plan        string
		expectError bool
		errorMsg    string
	}{
		{"Valid FreePlan", domain.FreePlan, false, ""},
		{"Valid BusinessPlan", domain.BusinessPlan, false, ""},
		{"Valid EnterprisePlan", domain.EnterprisePlan, false, ""},
		{"Invalid Plan", "UnknownPlan", true, "invalid user plan provided"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Plan(tt.plan)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSaveUser(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		plan        string
		expectError bool
		errorMsg    string
	}{
		{"Valid SaveUser", "user@example.com", "password123", domain.FreePlan, false, ""},
		{"Invalid email", "invalid-email", "password123", domain.FreePlan, true, "invalid email format"},
		{"Short password", "user@example.com", "short", domain.FreePlan, true, "password is too short"},
		{"Empty password", "user@example.com", "", domain.FreePlan, true, "empty password provided"},
		{"Invalid plan", "user@example.com", "password123", "UnknownPlan", true, "invalid user plan provided"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveUser(tt.email, tt.password, tt.plan)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
