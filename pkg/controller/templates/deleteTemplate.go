package templates

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	templateService "github.com/hngprojects/hng_boilerplate_golang_web/services/templates"
)

func (base *Controller) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")

	if !utility.IsValidUUID(id) {
		base.Logger.Error("Invalid id")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid id", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err := templateService.DeleteTemplate(base.Db.Postgresql, id)
	if err != nil {
		base.Logger.Error("Failed to delete template")
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to delete template", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Template deleted successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Template deleted successfully", nil)
	c.JSON(http.StatusOK, rd)
}
