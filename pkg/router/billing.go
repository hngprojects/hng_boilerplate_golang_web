package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/billing"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Billing(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	billing := billing.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	billingUrl := r.Group(fmt.Sprintf("%v", ApiVersion))

	{
		billingUrl.POST("/billing-plans", billing.CreateBilling)
		billingUrl.DELETE("/billing-plans/:id", billing.DeleteBilling)
		billingUrl.GET("/billing-plans", billing.GetBillings)
		billingUrl.GET("/billing-plans/:id", billing.GetBillingById)
		billingUrl.PATCH("/billing-plans/:id", billing.UpdateBillingById)
	}

	return r
}
