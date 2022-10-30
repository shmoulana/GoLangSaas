package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/shmoulana/Redios/internal/constant"
	"github.com/shmoulana/Redios/internal/service/logger"
)

type QueueService struct {
	Redis *redis.Client
}

type Queue struct {
	TypeQueue string `json:"type_queue"` // for example email, longjob
	Id        *int   `json:"id"`         // id longjob
	Value     string `json:"value"`      // value for email or meta longjob
	Status    string `json:"status"`
}

func (s QueueService) InsertQueue(ctx context.Context, data Queue) (*string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf(constant.RedisQueue, uuid.New())

	result := s.Redis.Set(ctx, key, string(b), time.Hour*24)
	if result.Err() != nil {
		return nil, err
	}

	return &key, err
}

func (s QueueService) AllQueues(ctx context.Context) ([]string, error) {
	var keys []string

	result := s.Redis.Keys(ctx, "queue-*")
	if result.Err() != nil && result.Err() != redis.Nil {
		return nil, nil
	}

	keys = result.Val()

	return keys, nil
}

func (s QueueService) UpdateStatus(ctx context.Context, key string, status string) error {
	var data Queue
	result := s.Redis.Get(ctx, key)
	if result.Err() != nil && result.Err() == redis.Nil {
		return result.Err()
	}

	dataRaw := result.Val()

	err := json.Unmarshal([]byte(dataRaw), &data)
	if err != nil {
		return err
	}

	data.Status = status

	res := s.Redis.Set(ctx, key, data, 24*time.Hour)
	if res.Err() != nil {
		return result.Err()
	}

	return nil
}

func (s QueueService) DeleteQueue(ctx context.Context, key string) error {
	result := s.Redis.Del(ctx, key)
	if result.Err() != nil && result.Err() != redis.Nil {
		return result.Err()
	}

	return nil
}

func (s QueueService) GetQueue(ctx context.Context, key string) (*Queue, error) {
	var data Queue

	result := s.Redis.Get(ctx, key)
	if result.Err() != nil && result.Err() == redis.Nil {
		return nil, result.Err()
	}

	dataRaw := result.Val()

	err := json.Unmarshal([]byte(dataRaw), &data)
	if err != nil {
		return nil, nil
	}

	return &data, nil
}

type WorkerPool struct {
	QueueService  QueueService
	LoggerService logger.LoggerService
}

func (wp WorkerPool) Run(ctx context.Context) error {
	keys, err := wp.QueueService.AllQueues(ctx)
	if err != nil {
		return err
	}

	fmt.Println(keys)

	return nil
}
