package config

import "os"

// Добавляем Redis конфигурацию
type redisConfig interface {
	RedisURL() string
	RedisPassword() string
}

type RedisConfig struct{}

func NewRedisConfig() (*RedisConfig, error) {
	return &RedisConfig{}, nil
}

func (c *RedisConfig) RedisURL() string {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		return "localhost:6379" // значение по умолчанию
	}
	return url
}

func (c *RedisConfig) RedisPassword() string {
	return os.Getenv("REDIS_PASSWORD") // может быть пустым
}
