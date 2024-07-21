package blog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
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

func (base *Controller) Post(c *gin.Context) {
	var (
		blogRequest models.CreateBlogRequest
	)

	if err := c.ShouldBindJSON(&blogRequest); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", "Bad Request", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&blogRequest); err != nil {
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

	respData, err := service.CreateBlog(blogRequest, base.Db.Postgresql, userId)

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("Blog created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Blog created successfully", respData)

	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) Delete(c *gin.Context) {
	blogID := c.Param("id")

	if err := service.DeleteBlog(blogID, base.Db.Postgresql); err != nil {
		if err.Error() == "blog not found" {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Blog with the given Id does not exist", "Not Found", nil)
			c.JSON(http.StatusNotFound, rd)
			return
		}
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to deleted blog", "internal server error", nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Blog successfully deleted")
	rd := utility.BuildSuccessResponse(http.StatusAccepted, "Blog successfully deleted", "Accepted")

	c.JSON(http.StatusAccepted, rd)

}
