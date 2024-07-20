package jobpost

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/go-playground/validator/v10"
	jobpost "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/job_post"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        		*storage.Database
	Validator 		*validator.Validate
	Logger 			*utility.Logger
	ExtReq 			request.ExternalRequest
}

func (base *Controller) CreateJobPost(c *gin.Context) {
    //    jobpost, err := jobpost.CreateJobPost(c, base.Db.Postgresql)
}

func (base *Controller) FetchAllJobPost(c *gin.Context) {
    //    jobpost, err := jobpost.FetchAllJobPost(c, base.Db.Postgresql)
}

func (base *Controller) FetchJobPostById(c *gin.Context) {
	
	//    jobpost, err := jobpost.FetchJobPostById(c, base.Db.Postgresql, id)
}

func (base *Controller) UpdateJobPostById(c *gin.Context) {
	//    jobpost, err := jobpost.UpdateJobPostById(c, base.Db.Postgresql, id)
}

// func (base *Controller) DeleteJobPostById(c *gin.Context) {
// 	//    jobpost, err := jobpost.DeleteJobPostById(c, base.Db.Postgresql, id)
// }