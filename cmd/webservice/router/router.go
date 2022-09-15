package router

import (
	"github.com/gin-gonic/gin"
	"github.com/shmoulana/Redios/cmd/webservice/handler"
)

type Router struct {
	r *gin.Engine
}

func NewRouter(ginEngine *gin.Engine) Router {
	return Router{
		r: ginEngine,
	}
}

func (h *Router) InitRouter() {
	h.r.GET(PingPath, handler.PingHandler())
	h.r.POST(TenantPath, handler.CreateTenantHandler())
}
