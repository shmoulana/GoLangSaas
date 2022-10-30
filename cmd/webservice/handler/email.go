package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shmoulana/Redios/pkg/dto"
)

type EmailTestNowHandler func(ctx context.Context) error
type EmailTestQueueHandler func(ctx context.Context) error
type EmailNowHandler func(ctx context.Context) error
type EmailQueueHandler func(ctx context.Context) error

func HandlerEmailTestNow(handler EmailTestNowHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := handler(ctx)
		if err != nil {
			WriteErrorResponse(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, dto.BaseResponse{
			Data: "ok",
		})

		return
	}
}

func HandlerEmailTestQueue(handler EmailTestQueueHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := handler(ctx)
		if err != nil {
			WriteErrorResponse(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, dto.BaseResponse{
			Data: "ok",
		})

		return
	}
}

func HandlerEmailNow(handler EmailNowHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := handler(ctx)
		if err != nil {
			WriteErrorResponse(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, dto.BaseResponse{
			Data: "ok",
		})

		return
	}
}

func HandlerEmailQueue(handler EmailQueueHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := handler(ctx)
		if err != nil {
			WriteErrorResponse(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, dto.BaseResponse{
			Data: "ok",
		})

		return
	}
}
