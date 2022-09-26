package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shmoulana/Redios/pkg/database"
	"github.com/shmoulana/Redios/pkg/dto"
)

type CreateTenant func(ctx context.Context, payload dto.TenantRequestV1) error
type GetTenants func(ctx context.Context) ([]database.Tenant, error)
type UpdateTenant func(ctx context.Context, id string, payload dto.TenantRequestV1) (*string, error)
type GetTenantById func(ctx context.Context, id string) (*database.Tenant, error)
type DeleteTenant func(ctx context.Context, id string) (*string, error)

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

func GetTenantHandler(handler GetTenants) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenants, err := handler(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.BaseResponse{
			Data: tenants,
		})

		return
	}
}

func UpdateTenantHandler(handler UpdateTenant) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req dto.TenantRequestV1
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tenants, err := handler(c, id, req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.BaseResponse{
			Data: tenants,
		})

		return
	}
}

func DeleteTenantHandler(handler DeleteTenant) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		tenants, err := handler(c, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.BaseResponse{
			Data: tenants,
		})

		return
	}
}

func GetTenantByIdHandler(handler GetTenantById) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		tenants, err := handler(c, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.BaseResponse{
			Data: tenants,
		})

		return
	}
}
