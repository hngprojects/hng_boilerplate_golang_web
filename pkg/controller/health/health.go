// Package health provides health check endpoints for the application
package health

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/ping"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

// Controller handles health check related operations
type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

// Post godoc
// @Summary Perform a health check post operation
// @Description Receives a ping message and returns a success response if the ping is successful
// @Tags health
// @Accept json
// @Produce json
// @Param request body models.Ping true "Ping request"
// @Success 200 {object} utility.Response "Successful ping response"
// @Failure 400 {object} utility.Response "Bad request"
// @Failure 500 {object} utility.Response "Internal server error"
// @Router /health/post [post]
func (base *Controller) Post(c *gin.Context) {
	var (
		req = models.Ping{}
	)
	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	base.Logger.Info("ping successful")
	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successful", req.Message)
	c.JSON(http.StatusOK, rd)
}

// Get godoc
// @Summary Perform a health check get operation
// @Description Returns a success response if the ping is successful
// @Tags health
// @Produce json
// @Success 200 {object} utility.Response "Successful ping response"
// @Failure 500 {object} utility.Response "Internal server error"
// @Router /health/get [get]
func (base *Controller) Get(c *gin.Context) {
	if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	base.Logger.Info("ping successful")
	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successful", gin.H{"transactions": "transactions object"})
	c.JSON(http.StatusOK, rd)
}
