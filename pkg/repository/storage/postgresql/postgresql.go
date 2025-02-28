package postgresql

import (
	"gorm.io/gorm"
)

type Postgresql struct {
	Db *gorm.DB
}

func (p *Postgresql) DB() *gorm.DB {
	if p.Db == nil {
		panic("Database must be instantiated before use")
	}
	return p.Db
}
func NewPostgresqlConnection(db *gorm.DB) *Postgresql {
	return &Postgresql{Db: db}
}
