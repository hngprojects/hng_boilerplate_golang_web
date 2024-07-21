package main

import (
	// "fmt"
	// "log"

	// "github.com/go-playground/validator/v10"
	// "github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	// "github.com/hngprojects/hng_boilerplate_golang_web/internal/models/migrations"
	// "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	// "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	// "github.com/hngprojects/hng_boilerplate_golang_web/pkg/router"
	// "github.com/hngprojects/hng_boilerplate_golang_web/utility"

	//for issue test no db ready for normal run
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/auth"
)

func main() {
	// logger := utility.NewLogger() //Warning !!!!! Do not recreate this action anywhere on the app

	// configuration := config.Setup(logger, "./app")

	// postgresql.ConnectToDatabase(logger, configuration.Database)
	// validatorRef := validator.New()

	// db := storage.Connection()

	// if configuration.Database.Migrate {
	// 	migrations.RunAllMigrations(db)
	// }

	// r := router.Setup(logger, validatorRef, db, &configuration.App)

	// utility.LogAndPrint(logger, fmt.Sprintf("Server is starting at 127.0.0.1:%s", configuration.Server.Port))
	// log.Fatal(r.Run(":8000"))

	route := gin.Default()
	//OAuth implementation by BlacAc3
	route.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "HNGi Golang Boilerplate",
			"status":  http.StatusOK,
		})
	})
	route.GET("/api/v1/auth/login/google", auth.Handle_Google_Login)
	route.GET("/api/v1/auth/callback/google", auth.Handle_Google_Callback)
	route.POST("/api/v1/auth/token/refresh", auth.Handle_Token_Refresh)
	route.Run(":8000")

}
