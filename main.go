package main

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	docs "github.com/hngprojects/hng_boilerplate_golang_web/docs"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models/migrations"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models/seed"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/oauth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/router"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

// @title HNG Boilerplate Golang Web API
// @version 1.0
// @description This is a boilerplate for golang HNG Internship 11.0
// @schemes http https

func main() {
	docs.SwaggerInfo.Title = "HNG Boilerplate Golang Web API"
	docs.SwaggerInfo.Description = "This is a boilerplate for golang HNG Internship 11.0"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8019"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	logger := utility.NewLogger() //Warning !!!!! Do not recreate this action anywhere on the app
	configuration := config.Setup(logger, "./app")
	postgresql.ConnectToDatabase(logger, configuration.Database)
	validatorRef := validator.New()
	oauth.SetupOauth(logger, configuration.Oauth) // setup oauth
	db := storage.Connection()

	if configuration.Database.Migrate {
		migrations.RunAllMigrations(db)
		// call the seed function
		seed.SeedDatabase(db.Postgresql)
	}

	r := router.Setup(logger, validatorRef, db, &configuration.App)
	utility.LogAndPrint(logger, fmt.Sprintf("Server is starting at 127.0.0.1:%s", configuration.Server.Port))
	log.Fatal(r.Run(":" + configuration.Server.Port))
}
