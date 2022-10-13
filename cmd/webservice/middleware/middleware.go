package middleware

import (
	"github.com/go-redis/redis/v8"
	"github.com/shmoulana/Redios/pkg/utils/crypt"
)

type Middleware struct {
	CryptService crypt.CryptService
	Redis        *redis.Client
}
