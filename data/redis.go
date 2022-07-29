package data

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	redisPool *redis.Pool
)

func init() {
	var err error
	redisPool, err = newLocalRedis()
	//log.Println("redisPool is ", redisPool)
	if err != nil {
		logrus.Fatal("redis init fail, ", err)
		return
	}
}

// 本地单点redis
func newLocalRedis() (*redis.Pool, error) {
	pool := &redis.Pool{
		// Maximum number of connections allocated by the pool at a given time.
		// When zero, there is no limit on the number of connections in the pool.
		//最大活跃连接数，0代表无限
		MaxActive: 0,
		//最大闲置连接数
		// Maximum number of idle connections in the pool.
		MaxIdle: 2000,
		//闲置连接的超时时间
		// Close connections after remaining idle for this duration. If the value
		// is zero, then idle connections are not closed. Applications should set
		// the timeout to a value less than the server's timeout.
		IdleTimeout: time.Second * 100,
		//定义拨号获得连接的函数
		// Dial is an application supplied function for creating and configuring a
		// connection.
		//
		// The connection returned from Dial must not be in a special state
		// (subscribed to pubsub channel, transaction started, ...).
		Dial: func() (redis.Conn, error) {
			//fmt.Println("redis.Host is ", conf.GetRedis().Host)
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				fmt.Println("when connect redis, happen error: ", err)
				return nil, err
			}
			return c, err
		},
	}
	return pool, nil
}

func RedisPool() *redis.Pool {
	return redisPool
}
