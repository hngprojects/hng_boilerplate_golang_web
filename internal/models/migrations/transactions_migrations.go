package migrations

import "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		&models.NewsLetter{}, //
	}
}

func AlterColumnModels() []AlterColumn {
	return []AlterColumn{}
}
