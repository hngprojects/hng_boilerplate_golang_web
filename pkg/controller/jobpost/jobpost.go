package jobpost

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/jobpost"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
  Db        *storage.Database
  Validator *validator.Validate
  Logger    *utility.Logger
  ExtReq    request.ExternalRequest 
}

func (base *Controller) CreateJobPost(c *gin.Context) {
	var req models.JobPost

	if err := c.ShouldBindJSON(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	if err := service.ValidateCreateJobPost(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", err.Error(), err, nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	respData, err := service.CreateJobPost(req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to create job post", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Job post created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "", respData)
	c.JSON(http.StatusCreated, rd)
}

func (base *Controller) FetchAllJobPost(c *gin.Context) {

    respData, err := service.FetchAllJobPost(base.Db.Postgresql)

	if err != nil {
    rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch job posts", err, nil)
    c.JSON(http.StatusInternalServerError, rd)
    return
  }

  base.Logger.Info("posts fetched successfully")
  rd := utility.BuildSuccessResponse(http.StatusOK, "", respData)
  c.JSON(http.StatusOK, rd)
}

func (base *Controller) FetchJobPostByID(c *gin.Context) {
	id := c.Param("id") 
	if _, err := uuid.Parse(id); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid ID format", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	respData, err := service.FetchJobPostByID(base.Db.Postgresql, id)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch job post", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Job post fetched successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "", respData)
	c.JSON(http.StatusOK, rd)
}