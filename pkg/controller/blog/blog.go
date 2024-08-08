package blog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/blog"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) CreateBlog(c *gin.Context) {
	var (
		blogReq models.CreateBlogRequest
	)

	if err := c.ShouldBind(&blogReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&blogReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", "Bad Request", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	respData, err := service.CreateBlog(blogReq, base.Db.Postgresql, userId)

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("Blog created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "blog created successfully", respData)

	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) DeleteBlog(c *gin.Context) {
	blogID := c.Param("id")

	if _, err := uuid.Parse(blogID); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid blog id format", "failed to delete blog", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", "Bad Request", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	if err := service.DeleteBlog(blogID, userId, base.Db.Postgresql); err != nil {
		if err.Error() == "blog not found" {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), "failed to delete blog", nil)
			c.JSON(http.StatusNotFound, rd)
			return
		}
		if err.Error() == "user not authorised to delete blog" {
			rd := utility.BuildErrorResponse(http.StatusForbidden, "error", err.Error(), "failed to delete blog", nil)
			c.JSON(http.StatusNotFound, rd)
			return
		}
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "failed to delete blog", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("blog successfully deleted")
	rd := utility.BuildSuccessResponse(http.StatusNoContent, "", nil)
	c.JSON(http.StatusNoContent, rd)

}

func (base *Controller) GetBlogs(c *gin.Context) {
	blogs, paginationResponse, err := service.GetBlogs(base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "failed to fetch blogs", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	paginationData := map[string]interface{}{
		"current_page": paginationResponse.CurrentPage,
		"total_pages":  paginationResponse.TotalPagesCount,
		"page_size":    paginationResponse.PageCount,
		"total_items":  len(blogs),
	}

	base.Logger.Info("blogs retrieved successfully.")
	rd := utility.BuildSuccessResponse(http.StatusOK, "blogs retrieved successfully", blogs, paginationData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetBlogById(c *gin.Context) {
	blogID := c.Param("id")

	if _, err := uuid.Parse(blogID); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid blog id format", "failed to delete blog", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	blog, err := service.GetBlogById(blogID, base.Db.Postgresql)

	if err != nil {
		if err.Error() == "blog not found" {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), "failed to retrieve blog", nil)
			c.JSON(http.StatusNotFound, rd)
			return
		}
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "failed to retrieve blog", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("blog retrieved successfully.")
	rd := utility.BuildSuccessResponse(http.StatusOK, "blog retrieved successfully", blog)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) UpdateBlogById(c *gin.Context){
	blogID := c.Param("id")
	var req models.UpdateBlogRequest

	if _, err := uuid.Parse(blogID); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid blog id format", "failed to update blog", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := c.ShouldBind(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userID, err := middleware.GetUserClaims(c, base.Db.Postgresql, "user_id")
	if err != nil {
		if err.Error() == "user claims not found" {
			rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), "failed to update blog", nil)
			c.JSON(http.StatusNotFound, rd)
			return
		}
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), "failed to update blog", nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	userId := userID.(string)

	blog, err := service.UpdateBlogById(blogID, userId, req, base.Db.Postgresql)

	if err !=nil {
		if err.Error() == "blog not found" {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), "failed to update blog", nil)
			c.JSON(http.StatusNotFound, rd)
			return
		}
		if err.Error() == "user not authorised to update blog" {
			rd := utility.BuildErrorResponse(http.StatusForbidden, "error", err.Error(), "failed to update blog", nil)
			c.JSON(http.StatusNotFound, rd)
			return
		}
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "failed to update blog", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("blog updated successfully.")
	rd := utility.BuildSuccessResponse(http.StatusOK, "blog updated successfully", blog)
	c.JSON(http.StatusOK, rd)
}
