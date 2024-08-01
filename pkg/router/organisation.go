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
		organisationUrl.POST("/organizations", organisation.CreateOrganisation)
		organisationUrl.GET("/organizations/:org_id", organisation.GetOrganisation)
		organisationUrl.DELETE("/organizations/:org_id", organisation.DeleteOrganisation)
		organisationUrl.PATCH("/organizations/:org_id", organisation.UpdateOrganisation)
		organisationUrl.GET("/organizations/:org_id/users", organisation.GetUsersInOrganisation)
		organisationUrl.POST("/organizations/:org_id/roles", organisation.CreateOrgRole)
		organisationUrl.GET("/organizations/:org_id/roles", organisation.GetOrgRoles)
		organisationUrl.GET("/organizations/:org_id/roles/:role_id", organisation.GetAOrgRole)
		organisationUrl.DELETE("/organizations/:org_id/roles/:role_id", organisation.DeleteOrgRole)
		organisationUrl.PATCH("/organizations/:org_id/roles/:role_id", organisation.UpdateOrgRole)
		organisationUrl.PATCH("/organizations/:org_id/roles/:role_id/permissions", organisation.UpdateOrgPermissions)
	}

	organisationUrlSec := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin))

	{
		organisationUrlSec.POST("/organizations/:org_id/users", organisation.AddUserToOrganisation)
	}
	return r
}
