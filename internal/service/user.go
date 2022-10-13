package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shmoulana/Redios/internal/constant"
	"github.com/shmoulana/Redios/internal/model"
	"github.com/shmoulana/Redios/internal/repository"
	"github.com/shmoulana/Redios/internal/service/logger"
	"github.com/shmoulana/Redios/pkg/dto"
	"github.com/shmoulana/Redios/pkg/errors"
	"github.com/shmoulana/Redios/pkg/utils/crypt"
)

type UserService struct {
	UserRepository repository.UserRepository
	CryptService   crypt.CryptService
	Redis          *redis.Client
	LoggerService  logger.LoggerService
}

func (s UserService) Name() string {
	return "UserService"
}

func (s UserService) Type() string {
	return "Service"
}

func (s UserService) Create(ctx context.Context, payload dto.SignUpPayload) (*int, error) {
	result, err := s.CryptService.CreateSignPSS(payload.Password)
	if err != nil {
		s.LoggerService.Error(s, err, "Failed to creating sign PSS")
		return nil, err
	}

	newUser := model.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: result,
	}

	lastInsertedId, err := s.UserRepository.Create(ctx, newUser)
	if err != nil {
		s.LoggerService.Error(s, err, "Failed to creating user")
		return nil, err
	}

	return lastInsertedId, nil
}

func (s UserService) SignIn(ctx context.Context, payload dto.SignInPayload) (*dto.TokenResponse, error) {
	user, err := s.UserRepository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		s.LoggerService.Error(s, err, fmt.Sprintf("Failed to find user by email:%s", payload.Email))
		return nil, err
	}

	if user == nil {
		err = errors.ErrDataNotFound
		return nil, err
	}

	isSame, err := s.CryptService.Verify(payload.Password, user.Password)
	if err != nil {
		err = errors.ErrInvalidUsernameOrPassword
		return nil, err
	}

	if !isSame {
		err = errors.ErrInvalidUsernameOrPassword
		return nil, err
	}

	userByte, err := json.Marshal(user)
	if err != nil {
		s.LoggerService.Error(s, err, "Failed to converting to byte")
		return nil, err
	}

	token, err := s.CryptService.CreateJWTToken(time.Hour*time.Duration(1), string(userByte))
	if err != nil {
		s.LoggerService.Error(s, err, "Failed to creating JWT token")
		return nil, err
	}

	key := fmt.Sprintf(constant.RedisTokenKey, user.ID)

	_, err = s.Redis.Get(ctx, key).Result()
	if err != redis.Nil {
		_, err = s.Redis.Del(ctx, key).Result()
		if err != nil {
			s.LoggerService.Error(s, err, fmt.Sprintf("Failed to deleting redis by key=%s", key))
			return nil, err
		}
	}

	_, err = s.Redis.Set(ctx, key, *token, time.Hour*time.Duration(1)).Result()
	if err != nil {
		s.LoggerService.Error(s, err, fmt.Sprintf("Failed to set redis by key=%s", key))

		return nil, err
	}

	return &dto.TokenResponse{
		Token: *token,
	}, nil
}
