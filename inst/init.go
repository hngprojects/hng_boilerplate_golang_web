package inst

import (
	"github.com/go-redis/redis/v8"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	red "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/redis"
	"gorm.io/gorm"
)

var (
	DB  *postgresql.Postgresql
	RDB *red.Redis
)

// initialize Postgresql instance
func InitDB(db *gorm.DB) *postgresql.Postgresql {
	if DB == nil {
		DB = postgresql.NewPostgresqlConnection(db)
	}
	return DB
}

// initialize Redis instance
func InitRed(rdb *redis.Client) *red.Redis {
	if RDB == nil {
		RDB = red.NewRedisConnection(rdb)
	}
	return RDB
}
