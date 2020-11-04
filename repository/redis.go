package repository

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/juju/errors"
)

type RedisRepository struct {
	rdb *redis.Client
}

const (
	redisAddr = "0.0.0.0"
	redisPort = 6379
)

var ctx = context.Background()

func NewRedisRepository() *RedisRepository {
	return &RedisRepository{}
}

func (r *RedisRepository) Connect(password string) error {
	r.rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisAddr, redisPort),
		Password: password, // no password set
	})

	_, err := r.rdb.Ping(ctx).Result()
	if err != nil {
		return errors.Wrap(err, errors.New("unable to connect to redis"))
	}

	return nil
}

func (r *RedisRepository) SetKeyValue(key, value string) error {
	err := r.rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		return errors.Wrap(err, errors.Errorf("unable to set key:value pair: %s:%s", key, value))
	}

	return nil
}

func (r *RedisRepository) GetValue(key string) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", errors.Wrap(err, errors.Errorf("unable to get key: %s", key))
	}

	return val, nil
}
