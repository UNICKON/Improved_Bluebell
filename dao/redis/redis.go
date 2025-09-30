package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

var rdb *redis.Client

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})
	if err := rdb.Ping().Err(); err != nil {
		zap.L().Error("ping DB failed, err:%v\n", zap.Error(err))
	}
	return
}

func Close() {
	_ = rdb.Close()
}

// 分布式锁实现
func Lock(key string, ttl time.Duration) (bool, error) {
	// 使用SET NX PX实现分布式锁
	ok, err := rdb.SetNX(key, "locked", ttl).Result()
	return ok, err
}

func Unlock(key string) error {
	_, err := rdb.Del(key).Result()
	return err
}
