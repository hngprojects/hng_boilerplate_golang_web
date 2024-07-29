package oauth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func SetupOauth(logger *utility.Logger, configOauth config.OauthFields) {

	goth.UseProviders(

		google.New(
			configOauth.GOOGLE_CLIENT_ID,
			configOauth.GOOGLE_CLIENT_SECRET,
			config.Config.App.Url+"/api/v1/auth/social/google/callback",
			"email",
			"profile",
		),

		facebook.New(
			configOauth.FACEBOOK_CLIENT_ID,
			configOauth.FACEBOOK_CLIENT_SECRET,
			config.Config.App.Url+"/api/v1/auth/social/facebook/callback",
			"email",
		),
	)

	maxAge := 86400 * 30 // 30 days
	isProd := false      // Set to true when serving over HTTPS

	store := sessions.NewCookieStore([]byte(configOauth.SESSION_SECRET))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd
	gothic.Store = store

	utility.LogAndPrint(logger, "initialized oauth configs...")
}
