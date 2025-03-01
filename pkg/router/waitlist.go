package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/waitlist"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	waitlistService "github.com/hngprojects/hng_boilerplate_golang_web/services/waitlist"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Waitlist(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {

	waitlistServices := waitlistService.NewWaitlistService(db.Postgresql.DB())
	controller := waitlist.Controller{DB: db, Validator: validator, Logger: logger, WaitlistService: waitlistServices}

	waitlistURL := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		waitlistURL.POST("/waitlist", controller.Create)
		waitlistURL.GET("/waitlist",
			middleware.Authorize(db.Postgresql.DB(), models.RoleIdentity.SuperAdmin), controller.GetWaitLists)
	}
	return r
}
