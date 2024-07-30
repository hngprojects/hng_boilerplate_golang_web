package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/jobpost"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func JobPost(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	controller := jobpost.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}
	jobPostUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		jobPostUrl.POST("/jobs", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), controller.CreateJobPost)
		jobPostUrl.GET("/jobs", controller.FetchAllJobPost)
		jobPostUrl.GET("/jobs/:job_id", controller.FetchJobPostByID)
		jobPostUrl.PATCH("/jobs/:job_id", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), controller.UpdateJobPostByID)
		jobPostUrl.DELETE("/jobs/:job_id", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), controller.DeleteJobPostByID)
	}
	return r
}
