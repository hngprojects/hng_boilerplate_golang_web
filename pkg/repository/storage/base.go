package storage

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type DbConnection interface {
	NewDatabaseConnection(db *gorm.DB, logger *utility.Logger, config *config.Database) *Database
}

type Database struct {
	Postgresql database.DatabaseManager
	Redis      database.CacheManager
}

var DB *Database = &Database{}

func Connection() *Database {
	return DB
}
