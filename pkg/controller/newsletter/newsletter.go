package newsletter

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/newsletter"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"net/http"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

// GetNewsLetters godoc
// @Summary Get all newsletters
// @Description Retrieve all newsletter subscriptions
// @Tags newsletters
// @Accept json
// @Produce json
// @Success 200 {object} utility.Response
// @Failure 400,404,500 {object} utility.Response
// @Router /newsletters [get]
func (base *Controller) GetNewsLetters(c *gin.Context) {
	newslettersData, paginationResponse, code, err := service.GetNewsletters(c, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "Newsletters email retrieved successfully", newslettersData, paginationResponse)
	c.JSON(http.StatusOK, rd)
}

// DeleteNewsLetter godoc
// @Summary Delete a newsletter subscription
// @Description Delete a newsletter subscription by ID
// @Tags newsletters
// @Accept json
// @Produce json
// @Param id path string true "Newsletter ID"
// @Success 200 {object} utility.Response
// @Failure 400,404,500 {object} utility.Response
// @Router /newsletters/{id} [delete]
func (base *Controller) DeleteNewsLetter(c *gin.Context) {
	var (
		reqID = c.Param("id")
	)
	code, err := service.DeleteNewsLetter(reqID, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), nil, nil)
		c.JSON(code, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "Newsletter email deleted successfully", nil)
	c.JSON(http.StatusOK, rd)
}

// SubscribeNewsLetter godoc
// @Summary Subscribe to newsletter
// @Description Subscribe a new email to the newsletter
// @Tags newsletters
// @Accept json
// @Produce json
// @Param newsletter body models.NewsLetter true "Newsletter subscription details"
// @Success 201 {object} utility.Response
// @Failure 400,409,422,500 {object} utility.Response
// @Router /newsletters/subscribe [post]
func (base *Controller) SubscribeNewsLetter(c *gin.Context) {
	var (
		req = models.NewsLetter{}
	)
	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed",
			utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}
	err = service.NewsLetterSubscribe(&req, base.Db.Postgresql)
	if err != nil {
		if err == models.ErrEmailAlreadySubscribed {
			rd := utility.BuildErrorResponse(http.StatusConflict, "error", "Email already subscribed", nil, nil)
			c.JSON(http.StatusConflict, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to subscribe", err, nil)
			c.JSON(http.StatusBadRequest, rd)
		}
		return
	}
	base.Logger.Info("subscribed successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "subscribed successfully", nil)
	c.JSON(http.StatusCreated, rd)
}
