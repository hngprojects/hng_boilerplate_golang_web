package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/superadmin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/superadmin"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func SuperAdmin(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}

	superAdminService := service.NewSuperAdminService(db.Postgresql.DB())
	superAdmin := superadmin.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq, SuperAdminService: superAdminService}

	superadminUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql.DB(), models.RoleIdentity.SuperAdmin))
	userUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql.DB()))
	{
		superadminUrl.POST("/regions", superAdmin.AddToRegion)
		superadminUrl.POST("/timezones", superAdmin.AddToTimeZone)
		superadminUrl.PATCH("/timezones/:id", superAdmin.UpdateTimeZone)
		superadminUrl.POST("/languages", superAdmin.AddToLanguage)
		userUrl.GET("/regions", superAdmin.GetRegion)
		userUrl.GET("/timezones", superAdmin.GetTimeZone)
		userUrl.GET("/languages", superAdmin.GetLanguage)
	}
	return r
}
