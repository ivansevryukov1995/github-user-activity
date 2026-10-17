package config

import (
	"fmt"

	"github.com/ivansevryukov1995/github-user-activity/internal/pkg/cache"
)

type Config struct {
	cache string
}

// InitConfig initialization config
func InitConfig() (*Config, error) {
	// 1. Cache
	name, err := cache.CacheFilePath()
	if err != nil {
		return nil, fmt.Errorf("failed init cache: %w", err)
	}

	return &Config{
		cache: name,
	}, nil
}

// Cache return name file for cache
func (c *Config) Cache() string {
	return c.cache
}
