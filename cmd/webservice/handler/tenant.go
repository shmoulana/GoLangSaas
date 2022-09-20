package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shmoulana/Redios/pkg/dto"
)

type CreateTenant func(ctx context.Context, payload dto.TenantRequestV1) error

func CreateTenantHandler(handler CreateTenant) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TenantRequestV1
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := handler(c, req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.BaseResponse{
			Data: "OK",
		})
		return
	}
}
