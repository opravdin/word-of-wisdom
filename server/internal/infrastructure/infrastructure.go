package infrastructure

import (
	"context"
	"fmt"

	"github.com/opravdin/word-of-wisdom/internal/configuration/env"
	"github.com/opravdin/word-of-wisdom/internal/logger"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Redis  *redis.Client
	Logger logger.Logger
}

func NewStorageConfiguration(config env.RedisConfig, log logger.Logger) (*Storage, error) {
	log.Debug("Initializing Redis connection", "host", config.Host, "port", config.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error("Failed to connect to Redis", "error", err)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Connected to Redis successfully", "host", config.Host, "port", config.Port)
	return &Storage{
		Redis:  rdb,
		Logger: log,
	}, nil
}

func (s *Storage) Close() error {
	if s.Redis != nil {
		s.Logger.Debug("Closing Redis connection")
		err := s.Redis.Close()
		if err != nil {
			s.Logger.Error("Error closing Redis connection", "error", err)
		}
		return err
	}
	return nil
}
