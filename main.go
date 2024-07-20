package main

import (
	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/auth"
)

func main() {
	// logger := utility.NewLogger() //Warning !!!!! Do not recreate this action anywhere on the app

	// configuration := config.Setup(logger, "./app")

	// postgresql.ConnectToDatabase(logger, configuration.Database)
	// // validatorRef := validator.New()

	// db := storage.Connection()

	// if configuration.Database.Migrate {
	// 	migrations.RunAllMigrations(db)
	// }

	// r := router.Setup(logger, validatorRef, db, &configuration.App)

	// utility.LogAndPrint(logger, fmt.Sprintf("Server is starting at 127.0.0.1:%s", configuration.Server.Port))
	// log.Fatal(r.Run(":8080"))

	//OAuth implementation by BlacAc3
	router := gin.Default()
	router.GET("/api/v1/auth/login/google", auth.Handle_Google_Login)
	router.GET("/api/v1/auth/callback/google", auth.Handle_Google_Callback)
	router.POST("/api/v1/auth/token/refresh", auth.Handle_Token_Refresh)
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Run(":8000")

}
