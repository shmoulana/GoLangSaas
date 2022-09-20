package internal

import (
	"github.com/shmoulana/Redios/configs"
	"github.com/shmoulana/Redios/internal/service"
	"github.com/shmoulana/Redios/pkg/database"
)

type Transport struct {
	tenantService *service.TenantService
	databaseRepo  *database.DatabaseRepo
}

func (t Transport) GetTenantService(conf configs.Config) service.TenantService {
	if t.tenantService == nil {
		tenantService := service.TenantService{
			DatabaseRepo: t.GetDatabaseRepo(conf),
		}

		t.tenantService = &tenantService
	}

	return *t.tenantService
}

func (t Transport) GetDatabaseRepo(conf configs.Config) database.DatabaseRepo {
	if t.databaseRepo == nil {
		// if conf.Driver == "sqlite3"{
		db := database.NewSQLLiteDriver()
		// }

		t.databaseRepo = &db
	}

	return *t.databaseRepo
}
