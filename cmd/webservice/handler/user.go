package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shmoulana/Redios/pkg/dto"
)

type Create func(ctx context.Context, payload dto.SignUpPayload) (*int, error)
type SignIn func(ctx context.Context, payload dto.SignInPayload) (*dto.TokenResponse, error)

func SignInHandler(handler SignIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SignInPayload
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lastInsertedId, err := handler(c, req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.BaseResponse{
			Data: lastInsertedId,
		})
		return
	}
}

func SignUpHandler(handler Create) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SignUpPayload
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := handler(c, req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.BaseResponse{
			Data: token,
		})
		return
	}
}
