package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/shmoulana/Redios/internal/model"
	"github.com/shmoulana/Redios/internal/repository"
	"github.com/shmoulana/Redios/pkg/dto"
	"github.com/shmoulana/Redios/pkg/errors"
	"github.com/shmoulana/Redios/pkg/util/crypt"
)

type UserService struct {
	UserRepository repository.UserRepository
	CryptService   crypt.CryptService
}

func (s UserService) Create(ctx context.Context, payload dto.SignUpPayload) (*int, error) {
	result, err := s.CryptService.CreateSignPSS(payload.Password)
	if err != nil {
		return nil, err
	}

	newUser := model.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: result,
	}

	lastInsertedId, err := s.UserRepository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return lastInsertedId, nil
}

func (s UserService) SignIn(ctx context.Context, payload dto.SignInPayload) (*dto.TokenResponse, error) {
	user, err := s.UserRepository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
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
		return nil, err
	}

	token, err := s.CryptService.CreateJWTToken(time.Hour*time.Duration(1), string(userByte))
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		Token: *token,
	}, nil
}
