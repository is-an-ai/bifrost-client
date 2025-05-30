package platform

import (
	"os"
	"runtime"
)

// StorageConfig represents platform-specific storage configuration
type StorageConfig struct {
	DirPerm  os.FileMode
	FilePerm os.FileMode
}

// GetStorageConfig returns platform-specific storage configuration
func GetStorageConfig() StorageConfig {
	switch runtime.GOOS {
	case "windows":
		return StorageConfig{
			DirPerm:  0666, // Windows: read/write for all
			FilePerm: 0666,
		}
	case "darwin", "linux":
		return StorageConfig{
			DirPerm:  0700, // Unix: owner only
			FilePerm: 0600,
		}
	default:
		return StorageConfig{
			DirPerm:  0700, // Default to Unix permissions
			FilePerm: 0600,
		}
	}
}
