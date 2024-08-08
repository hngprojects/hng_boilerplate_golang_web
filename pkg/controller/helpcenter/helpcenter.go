package helpcenter

import (
	"net/http"
	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/helpcenter"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) CreateHelpCenterTopic(c *gin.Context) {
	var req models.CreateHelpCenter

	if err := c.ShouldBindJSON(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	
    claims, exists := c.Get("userClaims")

	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", exists, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)

	userId := userClaims["user_id"].(string)


	user, code, err := user.GetUser(userId, base.Db.Postgresql)
    if err != nil {
        c.JSON(code, utility.BuildErrorResponse(code, "error", err.Error(), "Bad Request", nil))
        return
    }

	req.Author = user.Name
	req.Title = utility.CleanStringInput(req.Title)
    req.Content = utility.CleanStringInput(req.Content)
	
	if err := base.Validator.Struct(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Input validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	respData, err := service.CreateHelpCenterTopic(req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to add Topic", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Topic added successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Topic added successfully", respData)
	c.JSON(http.StatusCreated, rd)
}

func (base *Controller) FetchAllTopics(c *gin.Context) {
	topics, paginationResponse, err := service.GetPaginatedTopics(c, base.Db.Postgresql)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Topics not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch topics", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}
	paginationData := map[string]interface{}{
		"current_page": paginationResponse.CurrentPage,
		"total_pages":  paginationResponse.TotalPagesCount,
		"page_size":    paginationResponse.PageCount,
		"total_items":  len(topics),
	}
	base.Logger.Info("Topics retrieved successfully.")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Topics retrieved successfully.", topics, paginationData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) FetchTopicByID(c *gin.Context) {
		id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid ID format", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	respData, err := service.FetchTopicByID(base.Db.Postgresql, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Topic not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch topic", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("Topic retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Topic retrieved successfully", respData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) SearchHelpCenterTopics(c *gin.Context) {
	query := c.Query("title")
	if query == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Query parameter is required", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	topics, paginationResponse, err := service.SearchHelpCenterTopics(c, base.Db.Postgresql, query)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "No topics found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to search topics", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}
	paginationData := map[string]interface{}{
		"current_page": paginationResponse.CurrentPage,
		"total_pages":  paginationResponse.TotalPagesCount,
		"page_size":    paginationResponse.PageCount,
		"total_items":  len(topics),
	}
	base.Logger.Info("Topics retrieved successfully.")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Topics retrieved successfully.", topics, paginationData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) UpdateHelpCenterByID(c *gin.Context) {
	var req models.HelpCenter
	id := c.Param("id")

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
	claims, exists := c.Get("userClaims")

	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", exists, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)

	userId := userClaims["user_id"].(string)


	user, code, err := user.GetUser(userId, base.Db.Postgresql)
    if err != nil {
        c.JSON(code, utility.BuildErrorResponse(code, "error", err.Error(), "Bad Request", nil))
        return
    }

	req.Author = user.Name

	if err := base.Validator.Struct(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	result, err := service.UpdateTopic(base.Db.Postgresql, req, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Topic not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to update Topic", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("Topic updated successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Topic updated successfully", result)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) DeleteTopicByID(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid ID format", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err := service.DeleteTopicByID(base.Db.Postgresql, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Topic not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to delete Topic", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("Topic deleted successfully")
	rd := utility.BuildSuccessResponse(http.StatusNoContent, "", nil)
	c.JSON(http.StatusNoContent, rd)

}