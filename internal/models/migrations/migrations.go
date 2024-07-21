package migrations

import "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.Organisation{},
		models.Profile{},
		models.Product{},
		models.User{},
		models.Role{},
		models.Permission{},
	} // an array of db models, example: User{}
}

func AlterColumnModels() []AlterColumn {
	return []AlterColumn{}
}
