package templates

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	templateService "github.com/hngprojects/hng_boilerplate_golang_web/services/templates"

)

func (base *Controller) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")

	//validate uuid
	if !utility.IsValidUUID(id) {
		base.Logger.Error("Invalid id")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid id", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	var req models.EmailTemplate

	if err := c.ShouldBindJSON(&req); err != nil {
		base.Logger.Error("Failed to parse request body")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err := base.Validator.Struct(&req)
	if err != nil {
		base.Logger.Error("Validation failed")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}


	template, err := templateService.UpdateTemplate(base.Db.Postgresql, id, req)
	if err != nil {
		base.Logger.Error("Failed to update template")
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to update template", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Template updated successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "success", "Template updated successfully", template)
	c.JSON(http.StatusOK, rd)
}
