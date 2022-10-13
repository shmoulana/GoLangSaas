package redis_database

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/shmoulana/Redios/configs"
	"github.com/shmoulana/Redios/pkg/utils"
)

type RedisConfig struct {
	Hostname string
	Port     string
	Password string
	Database int
}

func NewRedisConnection(config configs.Config) *redis.Client {
	redisOpt := redis.Options{}
	redisOpt.Addr = fmt.Sprintf("%s:%s", config.RedisHostname, config.RedisPort)
	redisOpt.Password = config.RedisPassword
	redisOpt.DB = utils.StringToInt(config.RedisDatabase, 0) //default database

	client := redis.NewClient(&redisOpt)

	return client
}
