package webservice

import (
	"github.com/gin-gonic/gin"
	"github.com/shmoulana/Redios/cmd/webservice/router"
)

func StartServer() {
	r := gin.Default()

	route := router.NewRouter(r)

	route.InitRouter()

	r.Run(":8086")
	return
}
