package internal

import (
	"github.com/shmoulana/Redios/cmd/webservice/middleware"
	"github.com/shmoulana/Redios/configs"
	"github.com/shmoulana/Redios/internal/repository"
	"github.com/shmoulana/Redios/internal/service"
	"github.com/shmoulana/Redios/pkg/database"
	"github.com/shmoulana/Redios/pkg/util/crypt"
)

type Transport struct {
	tenantService     *service.TenantService
	userService       *service.UserService
	cryptService      *crypt.CryptService
	databaseRepo      *database.DatabaseRepo
	middlewareService *middleware.Middleware

	tenantRepository *repository.TenantRepository
	userRepository   *repository.UserRepository
}

// ---------------- Main
func (t Transport) GetDatabaseRepo(conf configs.Config) database.DatabaseRepo {
	if t.databaseRepo == nil {
		var db database.DatabaseRepo

		if conf.DBDriver == "sqlite3" {
			db = database.NewSQLLiteDriver()
		} else if conf.DBDriver == "postgres" {
			db = database.NewPostgreDriver(conf)
		}

		t.databaseRepo = &db
	}

	return *t.databaseRepo
}

func (t Transport) GetMiddleware(conf configs.Config) middleware.Middleware {
	if t.middlewareService == nil {
		middlewareService := middleware.Middleware{
			CryptService: t.GetCryptService(conf),
		}

		t.middlewareService = &middlewareService
	}

	return *t.middlewareService
}

// ---------------- Service
func (t Transport) GetTenantService(conf configs.Config) service.TenantService {
	if t.tenantService == nil {
		tenantService := service.TenantService{
			TenantRepository: t.GetTenantRepository(conf),
		}

		t.tenantService = &tenantService
	}

	return *t.tenantService
}

func (t Transport) GetUserService(conf configs.Config) service.UserService {
	if t.userService == nil {
		userService := service.UserService{
			UserRepository: t.GetUserRepository(conf),
			CryptService:   t.GetCryptService(conf),
		}

		t.userService = &userService
	}

	return *t.userService
}

func (t Transport) GetCryptService(conf configs.Config) crypt.CryptService {
	if t.cryptService == nil {
		cryptSvc := crypt.NewCryptService(crypt.Params{
			Conf: &conf,
		})

		t.cryptService = &cryptSvc
	}

	return *t.cryptService
}

// ---------------- Repository
func (t Transport) GetTenantRepository(conf configs.Config) repository.TenantRepository {
	if t.tenantRepository == nil {
		repo := repository.NewTenantRepository(t.GetDatabaseRepo(conf))

		t.tenantRepository = &repo
	}

	return *t.tenantRepository
}

func (t Transport) GetUserRepository(conf configs.Config) repository.UserRepository {
	if t.userRepository == nil {
		repo := repository.NewUserRepository(t.GetDatabaseRepo(conf))

		t.userRepository = &repo
	}

	return *t.userRepository
}
