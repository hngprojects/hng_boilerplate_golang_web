package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/helpcenter"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func HelpCenter(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	controller := helpcenter.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}
	helpCenterUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		helpCenterUrl.POST("/help-center/topics", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), controller.CreateHelpCenterTopic)
		helpCenterUrl.GET("/help-center/topics", controller.FetchAllTopics)
		helpCenterUrl.GET("/help-center/topics/:id", controller.FetchTopicByID)
		helpCenterUrl.GET("/help-center/topics/search", controller.SearchHelpCenterTopics)
		helpCenterUrl.PATCH("/help-center/topics/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), controller.UpdateHelpCenterByID)
		helpCenterUrl.DELETE("/help-center/topics/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin), controller.DeleteTopicByID)
	}
	return r
}