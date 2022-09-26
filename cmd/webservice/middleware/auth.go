package middleware

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
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

	c.Set(constant.UserKeyContext, user)

	c.Next()
}
