package infrastructure

import "github.com/redis/go-redis/v9"

var RedisClient = ConnectRedis()

func ConnectRedis() *redis.Client {
	store := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81", DB: 0})
	return store
}
