package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	controllers "github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/role"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
)

func Role(r *gin.Engine, ApiVersion string, db *storage.Database) *gin.Engine {
	roleController := controllers.NewRoleController(db.Postgresql)

	roleGroup := r.Group(fmt.Sprintf("/%s", ApiVersion), middleware.Authorize())
	{
		roleGroup.POST("/roles", roleController.CreateRole)
	}

	return r
}
