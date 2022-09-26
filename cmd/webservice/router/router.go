package router

import (
	"github.com/gin-gonic/gin"
	"github.com/shmoulana/Redios/cmd/webservice/handler"
	"github.com/shmoulana/Redios/cmd/webservice/middleware"
	"github.com/shmoulana/Redios/internal/service"
)

type Router struct {
	r             *gin.Engine
	tenantService service.TenantService
	userService   service.UserService
	middleware    middleware.Middleware
}

type NewRouterParams struct {
	R             *gin.Engine
	TenantService service.TenantService
	UserService   service.UserService
	Middleware    middleware.Middleware
}

func NewRouter(params NewRouterParams) Router {
	return Router{
		r:             params.R,
		tenantService: params.TenantService,
		userService:   params.UserService,
		middleware:    params.Middleware,
	}
}

func (h *Router) InitRouter() {
	h.r.GET(PingPath, handler.PingHandler())
	h.r.POST(UserSignInPath, handler.SignInHandler(h.userService.SignIn))
	h.r.POST(UserSignUpPath, handler.SignUpHandler(h.userService.Create))

	// Authorization group path
	authPath := h.r.Group("/", h.middleware.Authorization)
	// authPath := h.r.Group("/")

	// Tenants
	authPath.POST(TenantPath, handler.CreateTenantHandler(h.tenantService.CreateTenant))
	authPath.GET(TenantPath, handler.GetTenantHandler(h.tenantService.GetTenants))
	authPath.GET(TenantWithIdPath, handler.GetTenantByIdHandler(h.tenantService.GetTenantById))
	authPath.PUT(TenantWithIdPath, handler.UpdateTenantHandler(h.tenantService.UpdateTenant))
	authPath.DELETE(TenantWithIdPath, handler.DeleteTenantHandler(h.tenantService.DeleteTenant))
}
