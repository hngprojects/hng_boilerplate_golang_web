package router

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	jobpost "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/job_post"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func JobPost(r *gin.Engine, ApiVersion, string, validator *validator.Validate, db *storage.Database, logger *utility.Logger)*gin.Engine {
		extReq := request.ExternalRequest{Logger: logger, Test: false}
		result := jobpost.Controller(Db: db, validator: validator, logger:logger, extReq:extReq)
		jobPostUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
		{
			jobPostUrl.POST("/jobs", result.CreateJobPost)
			jobPostUrl.GET("/jobs", result.FetchAllJobPost)
			jobPostUrl.GET("/jobs/:id", result.FetchJobPostById)
			jobPostUrl.PATCH("/jobs/:id", result.UpdateJobPostById)
			// jobPostUrl.DELETE("/jobs/:id", result.DeleteJobPostById)
		}
		return r
}