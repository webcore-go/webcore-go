package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/semanggilab/webcore-go/app/config"
)

// Redis represents shared Redis connection
type Redis struct {
	Client *redis.Client
}

// NewRedis creates a new Redis connection
func NewRedis(config config.RedisConfig) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	return &Redis{Client: client}
}

// NewRedis creates a new Redis connection
func OpenRedis(config config.RedisConfig) (*Redis, error) {
	r := NewRedis(config)

	err := r.Connect()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Redis) Install(args ...any) error {
	// Tidak melakukan apa-apa
	return nil
}

func (r *Redis) Connect() error {
	// Test connection
	_, err := r.Client.Ping(r.Client.Context()).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}

// Close closes the Redis connection
func (r *Redis) Disconnect() error {
	return r.Client.Close()
}

func (r *Redis) Uninstall() error {
	// Tidak melakukan apa-apa
	return nil
}
