package middleware

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/shmoulana/Redios/cmd/webservice/handler"
	"github.com/shmoulana/Redios/internal/constant"
	"github.com/shmoulana/Redios/internal/model"
	"github.com/shmoulana/Redios/pkg/errors"
)

func (m Middleware) Authorization(c *gin.Context) {
	// Authorization
	var user model.User

	header := c.Request.Header

	if len(header["Authorization"]) < 1 {
		handler.WriteErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	authStr := header["Authorization"][0]
	tokenArry := strings.Split(authStr, " ")

	if len(tokenArry) < 2 {
		handler.WriteErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	token := tokenArry[1]

	userStr, err := m.CryptService.VerifyJWTToken(token)
	if err != nil {
		handler.WriteErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	err = json.Unmarshal([]byte(userStr.(string)), &user)
	if err != nil {
		handler.WriteErrorResponse(c, err)
		return
	}

	key := fmt.Sprintf(constant.RedisTokenKey, user.ID)

	redisToken, err := m.Redis.Get(c, key).Result()
	if err == redis.Nil {
		err = errors.ErrAuthTokenExpired

		handler.WriteErrorResponse(c, err)
		return
	}

	if redisToken != token {
		err = errors.ErrUnauthorized
		handler.WriteErrorResponse(c, err)
		return
	}

	c.Set(constant.UserKeyContext, user)

	c.Next()
}
