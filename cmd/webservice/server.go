package webservice

import (
	"github.com/gin-gonic/gin"
	"github.com/shmoulana/Redios/cmd/webservice/router"
	"github.com/shmoulana/Redios/configs"
	"github.com/shmoulana/Redios/internal"
)

func StartServer(conf configs.Config) {
	r := gin.Default()
	fac := internal.Transport{}

	route := router.NewRouter(router.NewRouterParams{
		R:             r,
		TenantService: fac.GetTenantService(conf),
		UserService:   fac.GetUserService(conf),
		Middleware:    fac.GetMiddleware(conf),
	})

	route.InitRouter()

	r.Run(":8086")
	return
}
