package jobpost

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/jobpost"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type Controller struct {
  Db        *storage.Database
  Validator *validator.Validate
  Logger    *utility.Logger
  ExtReq    request.ExternalRequest 
}

func (base *Controller) CreateJobPost(c *gin.Context) {
	var body struct {
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Location    string  `json:"location" binding:"required"`
		Salary      float64 `json:"salary"`
		JobType     string  `json:"job_type"`
		CompanyName string  `json:"company_name"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&body); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation errors", validationErrors.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	jobPost := models.JobPost{
		Title:       body.Title,
		Description: body.Description,
		Location:    body.Location,
		Salary:      body.Salary,
		JobType:     body.JobType,
		CompanyName: body.CompanyName,
	}

	post, err := jobpost.CreateJobPost(base.Db.Postgresql, jobPost)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to create job post", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("post created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "", post)
	c.JSON(http.StatusCreated, rd)
}

func (base *Controller) FetchAllJobPost(c *gin.Context) {
  posts, err := jobpost.FetchAllJobPost(base.Db.Postgresql)
  if err != nil {
    rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch job posts", err, nil)
    c.JSON(http.StatusInternalServerError, rd)
    return
  }

  base.Logger.Info("posts fetched successfully")
  rd := utility.BuildSuccessResponse(http.StatusOK, "", posts)
  c.JSON(http.StatusOK, rd)
}

func (base *Controller) FetchJobPostById(c *gin.Context) {
	jobPostID := c.Param("id")

  if _, err := uuid.Parse(jobPostID); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid job post ID", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	post, err := jobpost.FetchJobPostById(base.Db.Postgresql, jobPostID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Job post not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
			return
		}
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch job post", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("post fetched successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "", post)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) UpdateJobPostById(c *gin.Context) {
	jobPostID := c.Param("id")

  if _, err := uuid.Parse(jobPostID); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid job post ID", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}


	var body struct {
		Title       *string  `json:"title"` 
		Description *string  `json:"description"`
		Location    *string  `json:"location"`
		Salary      *float64 `json:"salary"`      
		JobType     *string  `json:"job_type"`
		CompanyName *string  `json:"company_name"`
	}

	if err := c.BindJSON(&body); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	updates := map[string]interface{}{}
	if body.Title != nil {
		updates["title"] = *body.Title
	}
	if body.Description != nil {
		updates["description"] = *body.Description
	}
	if body.Location != nil {
		updates["location"] = *body.Location
	}
	if body.Salary != nil {
		updates["salary"] = *body.Salary
	}
	if body.JobType != nil {
		updates["job_type"] = *body.JobType
	}
	if body.CompanyName != nil {
		updates["company_name"] = *body.CompanyName
	}

	post, err := jobpost.UpdateJobPostById(base.Db.Postgresql, jobPostID, updates)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Job post not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
			return
		}
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to update job post", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("post updated successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "", post)
	c.JSON(http.StatusOK, rd)
}