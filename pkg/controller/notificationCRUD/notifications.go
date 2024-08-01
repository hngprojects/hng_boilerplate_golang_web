package notificationCRUD

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	notificationcrud "github.com/hngprojects/hng_boilerplate_golang_web/services/notificationCRUD"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) UpdateNotificationSettings(c *gin.Context) {
	var req models.NotificationSettings

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	if err := c.ShouldBindJSON(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	_, err := notificationcrud.GetNotificationSettings(base.Db.Postgresql, userId)

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch notification settings", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	updated, err := notificationcrud.UpdateNotificationSettings(base.Db.Postgresql, req, userId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to update notification settings", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Notification settings updated successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Notification settings updated successfully", updated)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetNotificationSettings(c *gin.Context) {
	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	_, err := notificationcrud.GetNotificationSettings(base.Db.Postgresql, userId)
	if err != nil {
		notificationSettings := models.NotificationSettings{
			ID:                                   utility.GenerateUUID(),
			UserID:                               userId,
			MobilePushNotifications:              false,
			EmailNotificationActivityInWorkspace: false,
			EmailNotificationAlwaysSendEmailNotifications:  false,
			EmailNotificationEmailDigest:                   false,
			EmailNotificationAnnouncementAndUpdateEmails:   false,
			SlackNotificationsActivityOnYourWorkspace:      false,
			SlackNotificationsAlwaysSendEmailNotifications: false,
			SlackNotificationsAnnouncementAndUpdateEmails:  false,
		}

		_, err := notificationSettings.CreateNotificationSettings(base.Db.Postgresql)
		if err != nil {
			base.Logger.Info("Failed to create notification settings")
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to create notification settings", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
			return
		}
	}

	notificationSettings, err := notificationcrud.GetNotificationSettings(base.Db.Postgresql, userId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch notification settings", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Notification settings retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Notification settings retrieved successfully", notificationSettings)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) CreateNotification(c *gin.Context) {
	var req models.NotificationReq

	if err := c.ShouldBindJSON(&req); err != nil {
		base.Logger.Error("Invalid request body", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&req); err != nil {
		base.Logger.Error("Validation failed", err)
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	respData, err := notificationcrud.CreateNotification(base.Db.Postgresql, req, userId)
	base.Logger.Error("Failed to create notification", err)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to create notification", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	base.Logger.Info("Notification created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Notification created successfully", respData)
	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) FetchAllNotifications(c *gin.Context) {
	respData, addedData, err := notificationcrud.GetAllNotifications(c, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch notifications", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	data := map[string]interface{}{
		"total_notification_count":        addedData["total_count"],
		"total_unread_notification_count": addedData["unread_count"],
		"notifications":                   respData,
	}

	base.Logger.Info("Notifications retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Notifications retrieved successfully", data)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) FetchUnReadNotifications(c *gin.Context) {
	is_read := c.Param("is_read")

	if is_read != "true" && is_read != "false" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid Query Parameter", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	respData, addedData, err := notificationcrud.GetUnreadNotifications(c, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to fetch unread notifications", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	data := map[string]interface{}{
		"total_notification_count":        addedData["total_count"],
		"total_unread_notification_count": addedData["unread_count"],
		"notifications":                   respData,
	}

	base.Logger.Info("Unread Notifications retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Unread Notifications retrieved successfully", data)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) UpdateNotification(c *gin.Context) {
	var req models.Notification
	id := c.Param("notification_id")

	if _, err := uuid.Parse(id); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid ID format", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	result, err := notificationcrud.UpdateNotification(base.Db.Postgresql, req, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "notification not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to update notification", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("Notification updated successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Notification updated successfully", result)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) DeleteNotification(c *gin.Context) {
	id := c.Param("notification_id")

	if _, err := uuid.Parse(id); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid ID format", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err := notificationcrud.DeleteNotification(base.Db.Postgresql, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Notification not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to delete Notification", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("Notification deleted successfully")
	rd := utility.BuildSuccessResponse(http.StatusNoContent, "", nil)
	c.JSON(http.StatusNoContent, rd)

}
