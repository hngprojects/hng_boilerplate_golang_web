package migrations

import "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.Organisation{},
		models.User{},
		models.Profile{},
		models.Product{},
	} // an array of db models, example: User{}
}

func AlterColumnModels() []AlterColumn {
	return []AlterColumn{}
}
