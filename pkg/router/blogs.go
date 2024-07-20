package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/blogs"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Blog(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {

	blogs := blogs.Controller{Db: db, Validator: validator, Logger: logger}

	blogsUrl := r.Group(fmt.Sprintf("%v", ApiVersion))


	{
		blogsUrl.DELETE("/blogs/:id", blogs.Delete)
	}

	return r
}
