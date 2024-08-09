package migrations

import "github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.Testimonial{},
		models.SqueezeUser{},
		models.Blog{},
		models.AccessToken{},
		models.Role{},
		models.Organisation{},
		models.OrgRole{},
		models.Permission{},
		models.Profile{},
		models.Product{},
		models.User{},
		models.Invitation{},
		models.PasswordReset{},
		models.MagicLink{},
		models.WaitlistUser{},
		models.NewsLetter{},
		models.JobPost{},
		models.FAQ{},
		models.Language{},
		models.Timezone{},
		models.Region{},
		models.EmailTemplate{},
		models.UserRegionTimezoneLanguage{},
		models.Notification{},
		models.NotificationSettings{},
		models.HelpCenter{},
		models.ContactUs{},
		models.Billing{},
		models.DataPrivacySettings{},
		models.Key{},
	} // an array of db models, example: User{}
}

func AlterColumnModels() []AlterColumn {
	return []AlterColumn{}
}
