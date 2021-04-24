package persist

import (
	"github.com/Fakerr/sern/config"
	"log"

	"github.com/gomodule/redigo/redis"
)

// Redis connection pool
var pool *redis.Pool

func InitRedis() {
	pool = NewPool()
}

func NewPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			// If prod env, use DialUrl with the corresponding url
			if config.RedisURL != "" {

				log.Println("INFO: REDIS URI")

				c, err := redis.DialURL(config.RedisURL)
				if err != nil {
					panic(err.Error())
				}

				return c, err
			} else {
				c, err := redis.Dial("tcp", ":6379")
				if err != nil {
					panic(err.Error())
				}
				return c, err
			}

		},
	}
}

// Return a Redis connection Pool.
func GetRedisPool() *redis.Pool {
	return pool
}
