package blogs

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/blog"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) Delete(c *gin.Context) {
	blogID := c.Param("id")

	if err := blogs.DeleteBlog(blogID, base.Db.Postgresql); err != nil {
		if err.Error() == "blog not found" {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Blog with the given Id does not exist", fmt.Errorf("blog with the given Id does not exist"), nil)
			c.JSON(http.StatusNotFound, rd)
			return
		}
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Internal server error", fmt.Errorf("internal server error"), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Blog successfully deleted")
	rd := utility.BuildSuccessResponse(http.StatusAccepted, "Blog successfully deleted", "")

	c.JSON(http.StatusAccepted, rd)

}
