package main

import (
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/gin-gonic/gin"
=======
	// "fmt"
	// "log"

>>>>>>> feature/google_OAuth2
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
<<<<<<< HEAD
=======
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"

>>>>>>> c42f8ea65a2d0b943e9c08a5a2b5c20c9f1f1ad6
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models/migrations"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models/seed"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
=======
>>>>>>> feature/google_OAuth2
)

func main() {
	// logger := utility.NewLogger() //Warning !!!!! Do not recreate this action anywhere on the app

	// configuration := config.Setup(logger, "./app")

<<<<<<< HEAD
	postgresql.ConnectToDatabase(logger, configuration.Database)
<<<<<<< HEAD
=======
	// postgresql.ConnectToDatabase(logger, configuration.Database)
>>>>>>> feature/google_OAuth2
	// validatorRef := validator.New()
=======

	validatorRef := validator.New()
>>>>>>> c42f8ea65a2d0b943e9c08a5a2b5c20c9f1f1ad6

	// db := storage.Connection()

<<<<<<< HEAD
	if configuration.Database.Migrate {
		migrations.RunAllMigrations(db)

		// call the seed function
		seed.SeedDatabase(db.Postgresql)
	}
=======
	// if configuration.Database.Migrate {
	// 	migrations.RunAllMigrations(db)
	// }
>>>>>>> feature/google_OAuth2

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
