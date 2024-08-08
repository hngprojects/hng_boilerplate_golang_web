package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/notificationCRUD"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Notification(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	controller := notificationCRUD.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}
	notificationUrl := r.Group(fmt.Sprintf("%v/notifications", ApiVersion), middleware.Authorize(db.Postgresql))
	{
		notificationUrl.POST("/global", controller.CreateNotification)
		notificationUrl.GET("/all", controller.FetchAllNotifications)
		notificationUrl.GET("/unread", controller.FetchUnReadNotifications)
		notificationUrl.PATCH("/:notificationId", controller.UpdateNotification)
		notificationUrl.DELETE("/clear", controller.DeleteNotification)

	}
	return r
}

func NotificationSettings(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	controller := notificationCRUD.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}
	notificationUrl := r.Group(fmt.Sprintf("%v/settings", ApiVersion), middleware.Authorize(db.Postgresql))
	{
		notificationUrl.GET("/notification-settings", controller.GetNotificationSettings)
		notificationUrl.PATCH("/notification-settings", controller.UpdateNotificationSettings)
	}
	return r
}
