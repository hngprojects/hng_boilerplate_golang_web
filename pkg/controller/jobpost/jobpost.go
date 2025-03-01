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
	"gorm.io/gorm"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
	JobPostService service.JobPostService
}

func (base *Controller) CreateJobPost(c *gin.Context) {
	var req models.CreateJobPostModel

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

	jobPostService := service.NewJobPostService(base.Db.Postgresql.DB())

	respData, err := jobPostService.CreateJobPost(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to create job post", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Job created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Job created successfully", respData)
	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) FetchAllJobPost(c *gin.Context) {
	
	jobPostService := service.NewJobPostService(base.Db.Postgresql.DB())
	jobPosts, paginationResponse, err := jobPostService.GetPaginatedJobPosts(c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Jobs not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch job post", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}
	paginationData := map[string]interface{}{
		"current_page": paginationResponse.CurrentPage,
		"total_pages":  paginationResponse.TotalPagesCount,
		"page_size":    paginationResponse.PageCount,
		"total_items":  len(jobPosts),
	}
	base.Logger.Info("Job listings retrieved successfully.")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Job listings retrieved successfully.", jobPosts, paginationData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) FetchJobPostByID(c *gin.Context) {
	id := c.Param("job_id")
	if _, err := uuid.Parse(id); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid ID format", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	jobPostService := service.NewJobPostService(base.Db.Postgresql.DB())
	respData, err := jobPostService.FetchJobPostByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Job post not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch job post", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("Job post retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Job post retrieved successfully", respData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) UpdateJobPostByID(c *gin.Context) {
	var req models.JobPost
	id := c.Param("job_id")

	if _, err := uuid.Parse(id); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid ID format", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

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

	jobPostService := service.NewJobPostService(base.Db.Postgresql.DB())
	result, err := jobPostService.UpdateJobPost(req, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Job post not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to update job post", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("Job post updated successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Job post updated successfully", result)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) DeleteJobPostByID(c *gin.Context) {
	id := c.Param("job_id")

	if _, err := uuid.Parse(id); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid ID format", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	jobPostService := service.NewJobPostService(base.Db.Postgresql.DB())
	err := jobPostService.DeleteJobPostByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Job post not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to delete job post", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("Job post deleted successfully")
	rd := utility.BuildSuccessResponse(http.StatusNoContent, "", nil)
	c.JSON(http.StatusNoContent, rd)

}
