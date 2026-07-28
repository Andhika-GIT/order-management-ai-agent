package database

import (
	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func NewRedisClient(cfg *configs.RedisClientConfig) (*redis.Client, error) {

	opt, err := redis.ParseURL(cfg.Addr)

	if err != nil {
		return nil, err
	}
	return redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		// Explicitly disable maintenance notifications
		// This prevents the client from sending CLIENT MAINT_NOTIFICATIONS ON
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	}), nil
}
