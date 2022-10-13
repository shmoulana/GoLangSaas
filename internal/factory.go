package internal

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"github.com/shmoulana/Redios/cmd/webservice/middleware"
	"github.com/shmoulana/Redios/configs"
	"github.com/shmoulana/Redios/internal/repository"
	"github.com/shmoulana/Redios/internal/service"
	"github.com/shmoulana/Redios/internal/service/logger"
	"github.com/shmoulana/Redios/internal/service/queue"
	"github.com/shmoulana/Redios/pkg/database"
	"github.com/shmoulana/Redios/pkg/database/es"
	"github.com/shmoulana/Redios/pkg/database/redis_database"
	"github.com/shmoulana/Redios/pkg/utils/crypt"
)

type Transport struct {
	tenantService *service.TenantService
	userService   *service.UserService
	cryptService  *crypt.CryptService
	queueService  *queue.QueueService
	loggerService *logger.LoggerService

	databaseRepo      *database.DatabaseRepo
	middlewareService *middleware.Middleware

	tenantRepository *repository.TenantRepository
	userRepository   *repository.UserRepository

	redisClient   *redis.Client
	elasticClient *elasticsearch.Client
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

func (t Transport) GetRedis(conf configs.Config) *redis.Client {
	if t.redisClient == nil {
		redisClient := redis_database.NewRedisConnection(conf)

		t.redisClient = redisClient
	}

	return t.redisClient
}

func (t Transport) GetElasticsearch(conf configs.Config) *elasticsearch.Client {
	if t.elasticClient == nil {
		elasticClient, err := es.NewClient(conf)
		if err != nil {
			panic(err)
		}

		t.elasticClient = elasticClient
	}

	return t.elasticClient
}

// ---------------- Service
func (t Transport) GetTenantService(conf configs.Config) service.TenantService {
	if t.tenantService == nil {
		tenantService := service.TenantService{
			TenantRepository: t.GetTenantRepository(conf),
			LoggerService:    t.GetLoggerService(conf),
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
			LoggerService:  t.GetLoggerService(conf),
			Redis:          t.GetRedis(conf),
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

func (t Transport) GetQueueService(conf configs.Config) queue.QueueService {
	if t.queueService == nil {
		queueService := queue.QueueService{
			Redis: t.GetRedis(conf),
		}

		t.queueService = &queueService
	}

	return *t.queueService
}

func (t Transport) GetLoggerService(conf configs.Config) logger.LoggerService {
	if t.loggerService == nil {
		loggerService := logger.LoggerService{
			ElasticClient: t.GetElasticsearch(conf),
		}

		t.loggerService = &loggerService
	}

	return *t.loggerService
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
