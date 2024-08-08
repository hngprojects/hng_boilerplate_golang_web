package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/blog"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Blog(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	blogs := blog.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	blogsAdminUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin))
	blogsUrl := r.Group(fmt.Sprintf("%v", ApiVersion))

	{
		blogsAdminUrl.POST("/blogs", blogs.CreateBlog)
		blogsAdminUrl.DELETE("/blogs/:id", blogs.DeleteBlog)
		blogsAdminUrl.PATCH("/blogs/edit/:id", blogs.UpdateBlogById)
	}

	{
		blogsUrl.GET("/blogs", blogs.GetBlogs)
		blogsUrl.GET("/blogs/:id", blogs.GetBlogById)
	}

	return r
}
