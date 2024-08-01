package templates

import (
	"net/http"

	"github.com/gin-gonic/gin"
	templateService "github.com/hngprojects/hng_boilerplate_golang_web/services/templates"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func (base *Controller) GetTemplates(c *gin.Context) {

	templates, err := templateService.GetTemplates(base.Db.Postgresql)
	if err != nil {
		base.Logger.Error("Failed to retrieve templates")
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Templates Successfully retrieved")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Templates Successfully retrieved", templates)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetTemplate(c *gin.Context) {
	id := c.Param("id")

	//validate uuid
	if !utility.IsValidUUID(id) {
		base.Logger.Error("Invalid id")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid id", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	template, err := templateService.GetTemplate(base.Db.Postgresql, id)
	if err != nil {
		base.Logger.Error("Failed to retrieve template")
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Template Successfully retrieved")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Template Successfully retrieved", template)
	c.JSON(http.StatusOK, rd)


}