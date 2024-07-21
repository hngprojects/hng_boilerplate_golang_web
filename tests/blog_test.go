package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/router"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db := storage.Connection()
	validator := validator.New()
	logger := utility.NewLogger()
	r := gin.Default()

	apiVersion := "v1"
	router.Blog(r, apiVersion, validator, db, logger)
	return r
}


