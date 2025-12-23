package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/semanggilab/webcore-go/app/config"
	"github.com/semanggilab/webcore-go/app/helper"
	"github.com/semanggilab/webcore-go/app/logger"
)

// RedisPool manages Redis connection pools
type RedisPool struct {
	master *redis.Client
	slaves []*redis.Client
	config config.RedisConfig
	logger *logger.Logger
}

// NewRedisPool creates a new Redis connection pool
func NewRedisPool(config config.RedisConfig, logger *logger.Logger) (*RedisPool, error) {
	pool := &RedisPool{
		config: config,
		logger: logger,
	}

	// Connect to master Redis
	if err := pool.connectMaster(); err != nil {
		return nil, fmt.Errorf("failed to connect to master Redis: %v", err)
	}

	// Connect to slave Redis servers if configured
	if len(config.SlaveHosts) > 0 {
		if err := pool.connectSlaves(); err != nil {
			logger.Warn("Failed to connect to some slave Redis servers, continuing with master only")
		}
	}

	return pool, nil
}

// connectMaster connects to the master Redis server
func (p *RedisPool) connectMaster() error {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", p.config.Host, p.config.Port),
		Password:     p.config.Password,
		DB:           p.config.DB,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to ping master Redis: %v", err)
	}

	p.master = client
	p.logger.Info("Successfully connected to master Redis")
	return nil
}

// connectSlaves connects to slave Redis servers
func (p *RedisPool) connectSlaves() error {
	p.slaves = make([]*redis.Client, 0, len(p.config.SlaveHosts))

	for _, slaveConfig := range p.config.SlaveHosts {
		client := redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", slaveConfig.Host, slaveConfig.Port),
			Password:     slaveConfig.Password,
			DB:           slaveConfig.DB,
			MaxRetries:   3,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     10,
			MinIdleConns: 5,
		})

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := client.Ping(ctx).Result()
		if err != nil {
			p.logger.Warn(fmt.Sprintf("Failed to connect to slave Redis %s: %v", slaveConfig.Host, err))
			continue
		}

		p.slaves = append(p.slaves, client)
		p.logger.Info(fmt.Sprintf("Successfully connected to slave Redis: %s", slaveConfig.Host))
	}

	return nil
}

// GetClient returns a Redis client based on the operation type
func (p *RedisPool) GetClient() *redis.Client {
	// For read operations, use slaves if available
	// For write operations, use master
	// This is a simplified implementation
	return p.master
}

// GetMaster returns the master Redis client
func (p *RedisPool) GetMaster() *redis.Client {
	return p.master
}

// GetSlave returns a slave Redis client (round-robin)
func (p *RedisPool) GetSlave() *redis.Client {
	if len(p.slaves) == 0 {
		return p.master
	}
	// This is a simplified implementation - in production, you'd implement proper round-robin
	return p.slaves[0]
}

// Close closes all Redis connections
func (p *RedisPool) Disconnect() error {
	var errors []error

	// Close master connection
	if p.master != nil {
		if err := p.master.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close master connection: %v", err))
		}
	}

	// Close slave connections
	for i, slave := range p.slaves {
		if slave != nil {
			if err := slave.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close slave connection %d: %v", i, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors while closing Redis connections", len(errors))
	}

	return nil
}

// Health checks the health of all Redis connections
func (p *RedisPool) Health() map[string]any {
	health := make(map[string]any)

	// Check master
	if p.master != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		start := time.Now()
		_, err := p.master.Ping(ctx).Result()
		duration := time.Since(start)

		health["master"] = map[string]any{
			"status":    "healthy",
			"error":     helper.ErrToString(err),
			"latency":   duration.String(),
			"pool_size": p.master.PoolStats().TotalConns,
		}
	}

	// Check slaves
	slaveHealth := make([]map[string]any, 0, len(p.slaves))
	for i, slave := range p.slaves {
		if slave != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			start := time.Now()
			_, err := slave.Ping(ctx).Result()
			duration := time.Since(start)

			slaveHealth = append(slaveHealth, map[string]any{
				"id":        i,
				"host":      slave.Options().Addr,
				"status":    "healthy",
				"error":     helper.ErrToString(err),
				"latency":   duration.String(),
				"pool_size": slave.PoolStats().TotalConns,
			})
		}
	}

	health["slaves"] = slaveHealth
	health["total_slaves"] = len(slaveHealth)

	return health
}

// Set sets a key-value pair
func (p *RedisPool) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	return p.master.Set(ctx, key, value, expiration)
}

// Get gets a value by key
func (p *RedisPool) Get(ctx context.Context, key string) *redis.StringCmd {
	return p.master.Get(ctx, key)
}

// Delete deletes a key
func (p *RedisPool) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return p.master.Del(ctx, keys...)
}

// Exists checks if a key exists
func (p *RedisPool) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return p.master.Exists(ctx, keys...)
}

// Expire sets expiration for a key
func (p *RedisPool) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return p.master.Expire(ctx, key, expiration)
}

// TTL gets the time to live of a key
func (p *RedisPool) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return p.master.TTL(ctx, key)
}

// HSet sets a hash field
func (p *RedisPool) HSet(ctx context.Context, key string, field string, value any) *redis.IntCmd {
	return p.master.HSet(ctx, key, field, value)
}

// HGet gets a hash field
func (p *RedisPool) HGet(ctx context.Context, key string, field string) *redis.StringCmd {
	return p.master.HGet(ctx, key, field)
}

// HGetAll gets all fields and values of a hash
func (p *RedisPool) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	return p.master.HGetAll(ctx, key)
}

// HDel deletes a hash field
func (p *RedisPool) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	return p.master.HDel(ctx, key, fields...)
}

// LPush pushes a value to the left of a list
func (p *RedisPool) LPush(ctx context.Context, key string, values ...any) *redis.IntCmd {
	return p.master.LPush(ctx, key, values...)
}

// RPush pushes a value to the right of a list
func (p *RedisPool) RPush(ctx context.Context, key string, values ...any) *redis.IntCmd {
	return p.master.RPush(ctx, key, values...)
}

// LPop pops a value from the left of a list
func (p *RedisPool) LPop(ctx context.Context, key string) *redis.StringCmd {
	return p.master.LPop(ctx, key)
}

// RPop pops a value from the right of a list
func (p *RedisPool) RPop(ctx context.Context, key string) *redis.StringCmd {
	return p.master.RPop(ctx, key)
}

// LRange gets a range of values from a list
func (p *RedisPool) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return p.master.LRange(ctx, key, start, stop)
}

// SAdd adds a value to a set
func (p *RedisPool) SAdd(ctx context.Context, key string, members ...any) *redis.IntCmd {
	return p.master.SAdd(ctx, key, members...)
}

// SRem removes a value from a set
func (p *RedisPool) SRem(ctx context.Context, key string, members ...any) *redis.IntCmd {
	return p.master.SRem(ctx, key, members...)
}

// SMembers gets all members of a set
func (p *RedisPool) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return p.master.SMembers(ctx, key)
}

// ZAdd adds a value to a sorted set with score
func (p *RedisPool) ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return p.master.ZAdd(ctx, key, members...)
}

// ZRange gets a range of values from a sorted set by index
func (p *RedisPool) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return p.master.ZRange(ctx, key, start, stop)
}

// ZRangeByScore gets a range of values from a sorted set by score
func (p *RedisPool) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return p.master.ZRangeByScore(ctx, key, opt)
}

// Publish publishes a message to a channel
func (p *RedisPool) Publish(ctx context.Context, channel string, message any) *redis.IntCmd {
	return p.master.Publish(ctx, channel, message)
}

// Subscribe subscribes to a channel
func (p *RedisPool) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return p.master.Subscribe(ctx, channels...)
}

// PSubscribe subscribes to a pattern
func (p *RedisPool) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return p.master.PSubscribe(ctx, channels...)
}

// Incr increments a counter
func (p *RedisPool) Incr(ctx context.Context, key string) *redis.IntCmd {
	return p.master.Incr(ctx, key)
}

// IncrBy increments a counter by a value
func (p *RedisPool) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	return p.master.IncrBy(ctx, key, value)
}

// Decr decrements a counter
func (p *RedisPool) Decr(ctx context.Context, key string) *redis.IntCmd {
	return p.master.Decr(ctx, key)
}

// DecrBy decrements a counter by a value
func (p *RedisPool) DecrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	return p.master.DecrBy(ctx, key, value)
}
