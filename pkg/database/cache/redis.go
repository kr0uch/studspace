package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache: key not found")

type RedisConfig struct {
	Host     string `yaml:"redis_host" env:"REDIS_HOST" env-default:"localhost"`
	Port     string `yaml:"redis_port" env:"REDIS_PORT" env-default:"6379"`
	Password string `yaml:"redis_password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"redis_db" env:"REDIS_DB" env-default:"0"`

	MaxRetries   int           `yaml:"redis_max_retries" env:"REDIS_MAX_RETRIES" env-default:"3"`
	DialTimeout  time.Duration `yaml:"redis_dial_timeout" env:"REDIS_DIAL_TIMEOUT" env-default:"5s"`
	ReadTimeout  time.Duration `yaml:"redis_read_timeout" env:"REDIS_READ_TIMEOUT" env-default:"3s"`
	WriteTimeout time.Duration `yaml:"redis_write_timeout" env:"REDIS_WRITE_TIMEOUT" env-default:"3s"`
}

type RedisDB struct {
	client *redis.Client
}

func NewRedisDB(ctx context.Context, cfg RedisConfig) (*RedisDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	pingCtx, cancel := context.WithTimeout(ctx, cfg.DialTimeout)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis %s:%s/%d: %w", cfg.Host, cfg.Port, cfg.DB, err)
	}

	return &RedisDB{client: client}, nil
}

func (r *RedisDB) Close() error {
	return r.client.Close()
}

func (r *RedisDB) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisDB) SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	return r.client.SetNX(ctx, key, data, ttl).Result()
}

func (r *RedisDB) Get(ctx context.Context, key string, dest any) error {
	val, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return ErrCacheMiss
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(val, dest)
}

func (r *RedisDB) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}
