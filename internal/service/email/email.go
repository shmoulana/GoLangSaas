package email

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/shmoulana/Redios/configs"
	"github.com/shmoulana/Redios/internal/constant"
	"github.com/shmoulana/Redios/internal/repository"
	"github.com/shmoulana/Redios/internal/service/queue"
	"github.com/shmoulana/Redios/pkg/utils"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	TemplateRespository repository.TemplateRespository
	Redis               *redis.Client
	Config              configs.Config
	QueueService        queue.QueueService
}

type EmailJobValue struct {
	To  []string
	Msg string
}

// var auth smtp.Auth
var d *gomail.Dialer

func (s EmailService) SendNow(ctx context.Context, to []string, msg string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", "test@example.com")
	message.SetHeader("To", to...)
	message.SetHeader("Subject", "test subject")
	message.SetBody("text/html", msg)

	_, err := d.Dial()
	if err != nil {
		return err
	}

	err = d.DialAndSend(message)

	if err != nil {
		return err
	}

	return nil
}

func (s EmailService) Send(ctx context.Context, to []string, msg string) error {
	var canBeQueued bool = true

	redisStatus := s.Redis.Ping(ctx)
	if redisStatus.Err() != nil {
		canBeQueued = false
	}

	if canBeQueued {
		b, err := json.Marshal(EmailJobValue{
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

func (s EmailService) Init() {
	d = gomail.NewDialer(s.Config.EmailHost, utils.StringToInt(s.Config.EmailPort, 0), s.Config.EmailUsername, s.Config.EmailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
}
