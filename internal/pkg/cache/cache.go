package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ivansevryukov1995/github-user-activity/internal/model"
)

const CacheTTL = 60 * time.Second

func CacheFilePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed get directory cache: %w", err)
	}

	dir = filepath.Join(dir, "github-activity")

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed create directory cache: %w", err)
	}

	file := filepath.Join(dir, "cache.json")

	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return "", fmt.Errorf("failed create file cache: %w", err)
	}
	f.Close()

	return file, nil
}

func GetCache(nameFile, nameUser string) ([]model.Event, bool, error) {
	file, err := os.Open(nameFile)
	if err != nil {
		return nil, false, fmt.Errorf("failed open cache file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var foundUser *model.User

	for {
		var user model.User
		err := decoder.Decode(&user)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, false, fmt.Errorf("failed decode: %w", err)
		}

		if user.NameUser == nameUser {
			foundUser = &user
			if time.Since(user.CreatedAtCache) > CacheTTL {
				foundUser = nil
				continue
			}
			break
		}
	}

	if foundUser == nil {
		return nil, false, ErrUserNotFound
	}

	return foundUser.Events, true, nil
}

func SetCache(nameFile string, user *model.User) error {
	file, err := os.OpenFile(nameFile, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed open cache file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(user); err != nil {
		return fmt.Errorf("failed encode user to file: %w", err)
	}

	return nil
}

func DeleteCache() error {
	dir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("failed get directory cache: %w", err)
	}
	dir = filepath.Join(dir, "github-activity")

	file := filepath.Join(dir, "cache.json")

	if err := os.Remove(file); err != nil {
		return fmt.Errorf("failed clear cache: %w", err)
	}
	return nil
}
