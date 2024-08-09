package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Organisation(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	organisation := organisation.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	organisationUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User))
	{
		organisationUrl.POST("/organisations", organisation.CreateOrganisation)
		organisationUrl.GET("/organisations/:org_id", organisation.GetOrganisation)
		organisationUrl.DELETE("/organisations/:org_id", organisation.DeleteOrganisation)
		organisationUrl.PATCH("/organisations/:org_id", organisation.UpdateOrganisation)
		organisationUrl.GET("/organisations/:org_id/users", organisation.GetUsersInOrganisation)
		organisationUrl.POST("/organisations/:org_id/roles", organisation.CreateOrgRole)
		organisationUrl.GET("/organisations/:org_id/roles", organisation.GetOrgRoles)
		organisationUrl.GET("/organisations/:org_id/roles/:role_id", organisation.GetAOrgRole)
		organisationUrl.DELETE("/organisations/:org_id/roles/:role_id", organisation.DeleteOrgRole)
		organisationUrl.PATCH("/organisations/:org_id/roles/:role_id", organisation.UpdateOrgRole)
		organisationUrl.PATCH("/organisations/:org_id/roles/:role_id/permissions", organisation.UpdateOrgPermissions)
	}

	organisationUrlSec := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin))

	{
		organisationUrlSec.POST("/organisations/:org_id/users", organisation.AddUserToOrganisation)
	}
	return r
}
