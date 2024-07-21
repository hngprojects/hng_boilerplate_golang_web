package storage

import "gorm.io/gorm"
import (    
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
    "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)

type Database struct {
	Postgresql *gorm.DB
}

var DB *Database = &Database{}

func Connection() *Database {
	return DB
}


func (db *Database) CreateUserSubmission(submission *models.UserSubmission) error {
    return postgresql.CreateUserSubmission(db.Postgresql, submission)
}
