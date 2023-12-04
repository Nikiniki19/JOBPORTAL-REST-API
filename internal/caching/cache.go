package caching

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"job-portal-api/internal/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=cache.go -destination=cache_mock.go -package=caching
type Redis struct {
	rdb *redis.Client
}
type Cache interface {
	AddCache(ctx context.Context, jobid uint, jobData models.Job) error
	GetCache(ctx context.Context, jobid uint) (string, error)
	AddEmailToCache(ctx context.Context, email string, otp string) error
	GetEmailFromCache(ctx context.Context, otp string) (string, error)
}

func NewRedis(rdb *redis.Client) (Cache, error) {
	if rdb == nil {
		log.Info().Msg("redis cannot be nil")
		return nil, errors.New("---------------")
	}
	return &Redis{
		rdb: rdb,
	}, nil
}
func (re *Redis) AddCache(ctx context.Context, jobid uint, jobData models.Job) error {
	jobID := strconv.FormatUint(uint64(jobid), 10)
	val, err := json.Marshal(jobData)
	if err != nil {
		return err
	}
	err = re.rdb.Set(ctx, jobID, val, 15*time.Minute).Err()

	return err
}
func (re *Redis) GetCache(ctx context.Context, jobid uint) (string, error) {
	jobID := strconv.FormatUint(uint64(jobid), 10)
	str, err := re.rdb.Get(ctx, jobID).Result()
	return str, err
}

func (re *Redis) AddEmailToCache(ctx context.Context, email string, otp string) error {
	err := re.rdb.Set(ctx, email, otp, 5*time.Minute).Err()
	fmt.Println("[[[[[[[[[[]]]]]]]]]]", err)
	if err != nil {
		log.Err(err).Msg("error while adding==================")
		return fmt.Errorf("error while adding to redis : otp : %w = ", err)
	}
	return nil
}
func (re *Redis) GetEmailFromCache(ctx context.Context, email string) (string, error) {
	str, err := re.rdb.Get(ctx, email).Result()
	if err != nil {
		log.Err(err).Msg("error while getting=======================")
		return "", err
	}
	return str, nil
}
