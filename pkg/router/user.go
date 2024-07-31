package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func User(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	user := user.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	userUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User))
	adminUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin))
	{
		userUrl.GET("/users/:user_id", user.GetAUser)
		userUrl.DELETE("/users/:user_id", user.DeleteAUser)
		userUrl.PUT("/users/:user_id", user.UpdateAUser)
		userUrl.GET("/users/:user_id/organisations", user.GetAUserOrganisation)
		userUrl.PUT("/users/:user_id/roles/:role_id", user.AssignRoleToUser)
		userUrl.PUT("/users/:user_id/regions", user.UpdateUserRegion)
	}
	adminUrl.GET("/users", user.GetAllUsers)

	return r
}
