package notificationcrud

import (
	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func CreateNotification(db *gorm.DB, notificationReq models.NotificationReq, user_id string) (models.Notification, error) {
	// an instance of postgresql
	pdb := inst.InitDB(db)
	notification := models.Notification{
		ID:      utility.GenerateUUID(),
		UserID:  user_id,
		Message: notificationReq.Message,
	}

	createdNotification, err := notification.CreateNotification(pdb)
	if err != nil {
		return notification, err
	}

	return createdNotification, nil
}

func GetAllNotifications(c *gin.Context, db *gorm.DB) ([]models.Notification, map[string]int64, error) {
	var n models.Notification

	// an instance of postgresql
	pdb := inst.InitDB(db)

	notifications, addedData, err := n.FetchAllNotifications(pdb, c)
	if err != nil {
		return nil, nil, err
	}
	return notifications, addedData, nil
}

func GetUnreadNotifications(c *gin.Context, db *gorm.DB) ([]models.Notification, map[string]int64, error) {
	var notification models.Notification

	// an instance of postgresql
	pdb := inst.InitDB(db)

	notifications, addedData, err := notification.FetchUnReadNotifications(pdb, c)
	if err != nil {
		return nil, nil, err
	}
	return notifications, addedData, nil
}

func UpdateNotification(db *gorm.DB, notificationReq models.UpdateNotificationReq, ID string) (*models.Notification, error) {
	var n models.Notification

	// an instance of postgresql
	pdb := inst.InitDB(db)
	updatedNotification, err := n.UpdateNotification(pdb, notificationReq, ID)
	if err != nil {
		return updatedNotification, err
	}
	return updatedNotification, nil
}

func DeleteNotification(db *gorm.DB, ID string) error {
	var n models.Notification
	// an instance of postgresql
	pdb := inst.InitDB(db)

	return n.DeleteNotificationByUserID(pdb, ID)
}

func UpdateNotificationSettings(db *gorm.DB, notificationSettings models.NotificationSettings, ID string) (models.NotificationSettings, error) {

	var notificationSet models.NotificationSettings

	notificationSet.EmailNotificationActivityInWorkspace = notificationSettings.EmailNotificationActivityInWorkspace
	notificationSet.EmailNotificationAlwaysSendEmailNotifications = notificationSettings.EmailNotificationAlwaysSendEmailNotifications
	notificationSet.EmailNotificationAnnouncementAndUpdateEmails = notificationSettings.EmailNotificationAnnouncementAndUpdateEmails
	notificationSet.EmailNotificationEmailDigest = notificationSettings.EmailNotificationEmailDigest
	notificationSet.MobilePushNotifications = notificationSettings.MobilePushNotifications
	notificationSet.SlackNotificationsActivityOnYourWorkspace = notificationSettings.SlackNotificationsActivityOnYourWorkspace
	notificationSet.SlackNotificationsAlwaysSendEmailNotifications = notificationSettings.SlackNotificationsAlwaysSendEmailNotifications
	notificationSet.SlackNotificationsAnnouncementAndUpdateEmails = notificationSettings.SlackNotificationsAnnouncementAndUpdateEmails
	notificationSet.UserID = ID

	// an instance of postgresql
	pdb := inst.InitDB(db)
	updatedNotificationSettings, err := notificationSet.UpdateNotificationSettings(pdb, ID)
	if err != nil {
		return updatedNotificationSettings, err
	}
	return updatedNotificationSettings, nil
}

func GetNotificationSettings(db *gorm.DB, id string) (models.NotificationSettings, error) {
	var notificationSettings models.NotificationSettings

	// an instance of postgresql
	pdb := inst.InitDB(db)
	notificationSettings, err := notificationSettings.GetNotificationSettingsByID(pdb, id)
	if err != nil {
		return notificationSettings, err
	}
	return notificationSettings, nil
}
