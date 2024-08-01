package redis

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/go-redis/redis/v8"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

var Ctx = context.Background()

func ConnectToRedis(logger *utility.Logger, configDatabases config.Database) *redis.Client {
	dbsCV := configDatabases
	utility.LogAndPrint(logger, "connecting to redis server")
	connectedServer := connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.DB_NAME, dbsCV.DB_PORT, logger)

	utility.LogAndPrint(logger, "connected to redis server")

	storage.DB.Redis = connectedServer

	return connectedServer
}

func connectToDb(host, user, name, port string, logger *utility.Logger) *redis.Client {
	if _, err := strconv.Atoi(port); err != nil {
		u, err := url.Parse(port)
		if err != nil {
			utility.LogAndPrint(logger, fmt.Sprintf("parsing url %v to get port failed with: %v", port, err))
			panic(err)
		}

		detectedPort := u.Port()
		if detectedPort == "" {
			utility.LogAndPrint(logger, fmt.Sprintf("detecting port from url %v failed with: %v", port, err))
			panic(err)
		}
		port = detectedPort
	}
	db, err := strconv.Atoi(name)
	if err != nil {
		utility.LogAndPrint(logger, fmt.Sprintf("parsing url %v to get port failed with: %v", port, err))
		panic(err)
	}

	addr := fmt.Sprintf("%v:%v", host, port)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		Username: user,
		DB:       db,
	})

	if err := redisClient.Ping(Ctx).Err(); err != nil {
		utility.LogAndPrint(logger, fmt.Sprintln(addr))
		utility.LogAndPrint(logger, fmt.Sprintln("Redis db error: ", err))
	}

	pong, _ := redisClient.Ping(Ctx).Result()
	utility.LogAndPrint(logger, fmt.Sprintln("Redis says: ", pong))

	return redisClient
}
