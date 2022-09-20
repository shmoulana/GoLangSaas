package router

import (
	"github.com/gin-gonic/gin"
	"github.com/shmoulana/Redios/cmd/webservice/handler"
	"github.com/shmoulana/Redios/internal/service"
)

type Router struct {
	r             *gin.Engine
	tenantService service.TenantService
}

type NewRouterParams struct {
	R             *gin.Engine
	TenantService service.TenantService
}

func NewRouter(params NewRouterParams) Router {
	return Router{
		r:             params.R,
		tenantService: params.TenantService,
	}
}

func (h *Router) InitRouter() {
	h.r.GET(PingPath, handler.PingHandler())
	h.r.POST(TenantPath, handler.CreateTenantHandler(h.tenantService.CreateTenant))
}
