package migrations

import "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.Blog{},
		models.AccessToken{},
		models.Role{},
		models.Organisation{},
		models.Profile{},
		models.Product{},
		models.User{},
		models.Invitation{},
		models.PasswordReset{},
		models.MagicLink{},
		models.WaitlistUser{},
		models.NewsLetter{},
		models.JobPost{},
		models.PasswordReset{},
		models.MagicLink{},
	} // an array of db models, example: User{}
}

func AlterColumnModels() []AlterColumn {
	return []AlterColumn{}
}
