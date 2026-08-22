package helpers

import (
	"audituploader/internal/log"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func fileExists(file string) bool {
	cleanedPath := filepath.Clean(file)

	_, err := os.Stat(cleanedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		return false
	}

	return true
}

func FindFile(file string) (string, error) {
	if fileExists(file) {
		absPath, err := filepath.Abs(file)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path: %w", err)
		}
		log.Debug("Found file", "file", absPath)
		return absPath, nil
	}
	return "", fmt.Errorf("file does not exist: %s", file)
}
