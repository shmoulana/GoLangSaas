package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/goxiaoy/go-saas/data"
)

func CreateTenantHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sharedDsn := "./example.db"
		// connStrGen := saas.NewConnStrGenerator("./example-%s.db")

		conn := make(data.ConnStrings, 1)
		//default database
		conn.SetDefault(sharedDsn)

		// tenantStore := &saas.newtenants
		return
	}
}
