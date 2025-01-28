package validate

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFile(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		ownerID     int64
		expectError bool
		errorMsg    string
	}{
		{"ValidFile", 1, 1, false, ""},
		{"InvalidFileID", 0, 1, true, "invalid file id provided"},
		{"InvalidOwnerID", 1, 0, true, "invalid owner id provided 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := File(tt.id, tt.ownerID)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSaveFile(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		fileData    []byte
		expectError bool
		errorMsg    string
	}{
		{"ValidFile", "test.txt", []byte("file content"), false, ""},
		{"EmptyFilename", "", []byte("file content"), true, "empty filename is provided"},
		{"EmptyFileData", "test.txt", []byte(""), true, "empty file is provided"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveFile(tt.filename, tt.fileData)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
