package config

import (
	"fmt"
	"time"

	"github.com/ivansevryukov1995/github-user-activity/internal/pkg/cache"
)

type Config struct {
	cache *cache.Cache
}

// InitConfig initialization config
func InitConfig() (*Config, error) {
	// 1. Cache
	cache, err := cache.NewCache("github-activity", "cache.json", 120*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed init cache: %w", err)
	}

	return &Config{
		cache: cache,
	}, nil
}

// Cache return cache
func (c *Config) Cache() *cache.Cache {
	return c.cache
}
