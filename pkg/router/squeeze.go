package router

import (
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller"
    "github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
    "github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Squeeze(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) {
    squeeze := r.Group(ApiVersion + "/squeeze")
    {
        squeeze.POST("/", controller.HandleSqueeze)
    }
}