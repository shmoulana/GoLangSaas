package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shmoulana/Redios/pkg/dto"
)

func PingHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.BaseResponse{
			Data: "PONG!!!",
		})
	}
}
