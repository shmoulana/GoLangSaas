package email

import (
	"context"
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/go-redis/redis/v8"
	"github.com/shmoulana/Redios/configs"
	"github.com/shmoulana/Redios/internal/constant"
	"github.com/shmoulana/Redios/internal/repository"
	"github.com/shmoulana/Redios/internal/service/queue"
)

type EmailService struct {
	TemplateRespository repository.TemplateRespository
	Redis               *redis.Client
	Config              configs.Config
	QueueService        queue.QueueService
}

type queueValue struct {
	To  []string
	Msg string
}

var auth smtp.Auth

func (s EmailService) SendNow(ctx context.Context, to []string, msg []byte) error {
	auth = smtp.PlainAuth("", s.Config.EmailFrom, s.Config.EmailPassword, s.Config.EmailHost)
	addr := fmt.Sprintf("%s:%s", s.Config.EmailHost, s.Config.EmailPort)

	err := smtp.SendMail(addr, auth, s.Config.EmailFrom, to, msg)
	if err != nil {
		return err
	}

	return nil
}

func (s EmailService) Send(ctx context.Context, to []string, msg []byte) error {
	var canBeQueued bool = true

	redisStatus := s.Redis.Ping(ctx)
	if redisStatus.Err() != nil {
		canBeQueued = false
	}

	if canBeQueued {
		b, err := json.Marshal(queueValue{
			To:  to,
			Msg: string(msg),
		})

		if err != nil {
			return err
		}

		dataQueue := queue.Queue{
			TypeQueue: constant.TypeQueueEmail,
			Value:     string(b),
		}

		key, err := s.QueueService.InsertQueue(ctx, dataQueue)
		if err != nil {
			return err
		}

		fmt.Print(key)
	} else {
		err := s.Send(ctx, to, msg)
		if err != nil {
			return err
		}
	}

	return nil
}
