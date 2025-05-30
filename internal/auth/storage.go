package auth

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"bifrost-client/internal/platform"
)

// Storage defines the interface for token storage
type Storage interface {
	// SaveToken saves the authentication token
	SaveToken(ctx context.Context, token string) error

	// GetToken retrieves the authentication token
	GetToken(ctx context.Context) (string, error)

	// DeleteToken removes the authentication token
	DeleteToken(ctx context.Context) error
}

// LocalStorage implements the Storage interface using local storage
type LocalStorage struct {
	filePath string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage() (*LocalStorage, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Get platform-specific storage configuration
	config := platform.GetStorageConfig()

	// Create .bifrost directory if it doesn't exist
	bifrostDir := filepath.Join(homeDir, ".bifrost")
	if err := os.MkdirAll(bifrostDir, config.DirPerm); err != nil {
		return nil, err
	}

	return &LocalStorage{
		filePath: filepath.Join(bifrostDir, "auth.json"),
	}, nil
}

// SaveToken implements Storage interface
func (s *LocalStorage) SaveToken(ctx context.Context, token string) error {
	data := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	// Get platform-specific file permissions
	config := platform.GetStorageConfig()

	file, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, config.FilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(data)
}

// GetToken implements Storage interface
func (s *LocalStorage) GetToken(ctx context.Context) (string, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	var data struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return "", err
	}

	return data.Token, nil
}

// DeleteToken implements Storage interface
func (s *LocalStorage) DeleteToken(ctx context.Context) error {
	return os.Remove(s.filePath)
}
