package billing

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
	"github.com/hngprojects/hng_boilerplate_golang_web/services/billing"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) CreateBilling(c *gin.Context) {
	var (
		billingReq models.CreateBillingRequest
	)

	if err := c.ShouldBind(&billingReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&billingReq); err != nil {
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

	respData, err := billing.CreateBilling(billingReq, base.Db.Postgresql, userId)

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("billing created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "billing created successfully", respData)

	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) DeleteBilling(c *gin.Context) {
	billingID := c.Param("id")

	if _, err := uuid.Parse(billingID); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid billing id format", "failed to delete billing", nil)
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

	if err := billing.DeleteBilling(billingID, userId, base.Db.Postgresql); err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "billing not found", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	base.Logger.Info("billing successfully deleted")
	rd := utility.BuildSuccessResponse(http.StatusNoContent, "", nil)
	c.JSON(http.StatusNoContent, rd)

}

func (base *Controller) GetBillings(c *gin.Context) {
	billings_len, paginationResponse, err := billing.GetBillings(base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "failed to fetch billings", err, nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	paginationData := map[string]interface{}{
		"current_page": paginationResponse.CurrentPage,
		"total_pages":  paginationResponse.TotalPagesCount,
		"page_size":    paginationResponse.PageCount,
		"total_items":  billings_len,
	}

	base.Logger.Info("billings retrieved successfully.")
	rd := utility.BuildSuccessResponse(http.StatusOK, "billings retrieved successfully", billings_len, paginationData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetBillingById(c *gin.Context) {
	billingID := c.Param("id")

	if _, err := uuid.Parse(billingID); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid billing id format", "failed to delete billing", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	billing, err := billing.GetBillingById(billingID, base.Db.Postgresql)

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "billing not found", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	base.Logger.Info("billing retrieved successfully.")
	rd := utility.BuildSuccessResponse(http.StatusOK, "billing retrieved successfully", billing)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) UpdateBillingById(c *gin.Context) {
	billingID := c.Param("id")
	var req models.UpdateBillingRequest

	if _, err := uuid.Parse(billingID); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid billing id format", "failed to update billing", nil)
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
			rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), "failed to update billing", nil)
			c.JSON(http.StatusNotFound, rd)
			return
		}
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), "failed to update billing", nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	userId := userID.(string)

	billing, err := billing.UpdateBillingById(billingID, userId, req, base.Db.Postgresql)

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "billing not found", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	base.Logger.Info("billing updated successfully.")
	rd := utility.BuildSuccessResponse(http.StatusOK, "billing updated successfully", billing)
	c.JSON(http.StatusOK, rd)
}
