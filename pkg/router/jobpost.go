package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/jobpost"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func JobPost(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger)*gin.Engine {
		extReq := request.ExternalRequest{Logger: logger, Test: false}
		controller := jobpost.Controller{Db: db, Validator: validator, Logger:logger, ExtReq:extReq}
		jobPostUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
		{
			jobPostUrl.POST("/jobs", controller.CreateJobPost)
			jobPostUrl.GET("/jobs", controller.FetchAllJobPost)
			jobPostUrl.GET("/jobs/:id", controller.FetchJobPostByID)
		}
		return r
}