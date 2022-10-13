package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/shmoulana/Redios/internal/constant"
)

type QueueService struct {
	Redis *redis.Client
}

type Queue struct {
	TypeQueue string `json:"type_queue"`
	Id        *int   `json:"id"`
	Value     string `json:"value"`
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

func (s QueueService) queues(ctx context.Context) ([]string, error) {
	var keys []string

	result := s.Redis.Keys(ctx, "*")
	if result.Err() != nil && result.Err() != redis.Nil {
		return nil, nil
	}

	keys = result.Val()

	return keys, nil
}

func (s QueueService) Run(ctx context.Context) error {
	keys, err := s.queues(ctx)
	if err != nil {
		return err
	}

	fmt.Print(keys)
	return nil
}
