package migrations

import "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.AccessToken{},
		models.Role{},
		models.Organisation{},
		models.Profile{},
		models.Product{},
		models.User{},
		models.WaitlistUser{},
		models.NewsLetter{},
		models.JobPost{},
	} // an array of db models, example: User{}
}

func AlterColumnModels() []AlterColumn {
	return []AlterColumn{}
}
