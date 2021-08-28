package redis

import (
	"sync"

	"github.com/go-redis/redis/v8"

	"bigs-ci/config"
)

var Client *client


type client struct {
	Cli *redis.Client
}

func Init() {
	var once sync.Once
	once.Do(New)

}

func New()  {
	Client= &client{Cli: redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisConf.Addr,
		Password: config.Config.RedisConf.Password, // no password set
		DB:       config.Config.RedisConf.DB,       // use default DB
	})}

}
