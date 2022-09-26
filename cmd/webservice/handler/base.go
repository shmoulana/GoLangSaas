package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shmoulana/Redios/pkg/dto"
	"github.com/shmoulana/Redios/pkg/errors"
)

func WriteErrorResponse(c *gin.Context, er error) {
	errResp := errors.GetErrorResponse(er)
	resp := dto.BaseResponse{
		Error: &dto.ErrorResponse{
			Code:    errResp.Code,
			Message: errResp.Message,
		},
	}

	c.AbortWithStatusJSON(int(errResp.HTTPCode), resp)
	return
}
