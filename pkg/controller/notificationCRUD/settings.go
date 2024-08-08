package notificationCRUD

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	notificationcrud "github.com/hngprojects/hng_boilerplate_golang_web/services/notificationCRUD"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

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
