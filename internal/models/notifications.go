package models

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"gorm.io/gorm"
)

type Notification struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;unique;not null"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type NotificationSettings struct {
	ID                                             string `json:"id" gorm:"type:uuid;primaryKey;unique;not null"`
	UserID                                         string `json:"user_id" gorm:"type:uuid;not null"`
	MobilePushNotifications                        bool   `json:"mobile_push_notifications" gorm:"default:false"`
	EmailNotificationActivityInWorkspace           bool   `json:"email_notification_activity_in_workspace" gorm:"default:false"`
	EmailNotificationAlwaysSendEmailNotifications  bool   `json:"email_notification_always_send_email_notifications" gorm:"default:false"`
	EmailNotificationEmailDigest                   bool   `json:"email_notification_email_digest" gorm:"default:false"`
	EmailNotificationAnnouncementAndUpdateEmails   bool   `json:"email_notification_announcement_and_update_emails" gorm:"default:false"`
	SlackNotificationsActivityOnYourWorkspace      bool   `json:"slack_notifications_activity_on_your_workspace" gorm:"default:false"`
	SlackNotificationsAlwaysSendEmailNotifications bool   `json:"slack_notifications_always_send_email_notifications" gorm:"default:false"`
	SlackNotificationsAnnouncementAndUpdateEmails  bool   `json:"slack_notifications_announcement_and_update_emails" gorm:"default:false"`
}

type NotificationReq struct {
	Message string `json:"message"`
}

type UpdateNotificationReq struct {
	IsRead bool `json:"is_read"`
}

func (n *Notification) CreateNotification(db database.DatabaseManager) (Notification, error) {

	err := db.CreateOneRecord(&n)

	notification := Notification{
		ID:      n.ID,
		UserID:  n.UserID,
		Message: n.Message,
	}

	if err != nil {
		return notification, err
	}
	return notification, nil
}

func (n *Notification) GetNotificationByID(db database.DatabaseManager, ID string) (Notification, error) {
	var notification Notification

	err, _ := db.SelectOneFromDb(&notification, "id = ?", ID)
	if err != nil {
		return notification, err
	}
	return notification, nil
}

func (n *Notification) FetchAllNotifications(db database.DatabaseManager, c *gin.Context) ([]Notification, map[string]int64, error) {
	var notifications []Notification
	type additionalData map[string]int64

	err := db.SelectAllFromDb(nil, "", &notifications, "")
	if err != nil {
		return nil, additionalData{}, err
	}

	totalCount, err := db.CountRecords(&notifications)
	if err != nil {
		return nil, additionalData{}, err
	}

	unreadCount, err := db.CountSpecificRecords(&notifications, "is_read = false")
	if err != nil {
		return nil, additionalData{}, err
	}

	data := additionalData{
		"total_count":  totalCount,
		"unread_count": unreadCount,
	}

	return notifications, data, nil
}

func (n *Notification) FetchUnReadNotifications(db database.DatabaseManager, c *gin.Context) ([]Notification, map[string]int64, error) {
	var notifications []Notification
	type additionalData map[string]int64

	totalCount, err := db.CountRecords(&n)
	if err != nil {
		return nil, additionalData{}, err
	}

	unreadCount, err := db.CountSpecificRecords(&n, "is_read = false")
	if err != nil {
		return nil, additionalData{}, err
	}

	data := additionalData{
		"total_count":  totalCount,
		"unread_count": unreadCount,
	}

	err = db.SelectAllFromDb(nil, "", &notifications, "is_read = ?", false)
	if err != nil {
		return nil, additionalData{}, err
	}
	return notifications, data, nil
}

func (n *Notification) UpdateNotification(db database.DatabaseManager, notifReq UpdateNotificationReq, ID string) (*Notification, error) {

	exists := db.CheckExists(&n, "id = ?", ID)
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	oldNotification, err := n.GetNotificationByID(db, ID)
	if err != nil {
		return nil, err
	}

	res, err := db.UpdateFields(&n, notifReq, ID)
	if err != nil {
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, errors.New("failed to update notification")
	}

	n.ID = oldNotification.ID
	n.UserID = oldNotification.UserID
	n.Message = oldNotification.Message

	return n, nil
}

func (n *Notification) DeleteNotificationByUserID(db database.DatabaseManager, ID string) error {
	var notifications []Notification

	exists := db.CheckExists(&notifications, "user_id = ?", ID)
	if !exists {
		return gorm.ErrRecordNotFound
	}

	err := db.SelectAllFromDb(nil, "", &notifications, "user_id = ?", ID)
	if err != nil {
		return err
	}

	err = db.DeleteRecordFromDb(&notifications)
	if err != nil {
		return err
	}
	return nil
}

func (n *NotificationSettings) GetNotificationSettingsByID(db database.DatabaseManager, ID string) (NotificationSettings, error) {
	var notificationSettings NotificationSettings

	err, nerr := db.SelectOneFromDb(&notificationSettings, "user_id = ?", ID)
	if err != nil {
		return notificationSettings, nerr
	}
	return notificationSettings, nil
}

func (n *NotificationSettings) CreateNotificationSettings(db database.DatabaseManager) (NotificationSettings, error) {
	err := db.CreateOneRecord(&n)

	notificationSettings := NotificationSettings{
		ID:                                   n.ID,
		UserID:                               n.UserID,
		MobilePushNotifications:              n.MobilePushNotifications,
		EmailNotificationActivityInWorkspace: n.EmailNotificationActivityInWorkspace,
		EmailNotificationAlwaysSendEmailNotifications:  n.EmailNotificationAlwaysSendEmailNotifications,
		EmailNotificationEmailDigest:                   n.EmailNotificationEmailDigest,
		EmailNotificationAnnouncementAndUpdateEmails:   n.EmailNotificationAnnouncementAndUpdateEmails,
		SlackNotificationsActivityOnYourWorkspace:      n.SlackNotificationsActivityOnYourWorkspace,
		SlackNotificationsAlwaysSendEmailNotifications: n.SlackNotificationsAlwaysSendEmailNotifications,
		SlackNotificationsAnnouncementAndUpdateEmails:  n.SlackNotificationsAnnouncementAndUpdateEmails,
	}

	if err != nil {
		return notificationSettings, err
	}
	return notificationSettings, nil
}

func (n *NotificationSettings) UpdateNotificationSettings(db database.DatabaseManager, ID string) (NotificationSettings, error) {
	n.ID = ID

	exists := db.CheckExists(&NotificationSettings{}, "user_id = ?", ID)
	if !exists {
		return NotificationSettings{}, gorm.ErrRecordNotFound
	}

	_, err := db.SaveAllFields(n)
	if err != nil {
		return NotificationSettings{}, err
	}

	updatedNotification := NotificationSettings{}
	err = db.DB().First(&updatedNotification, "id = ?", ID).Error
	if err != nil {
		return NotificationSettings{}, err
	}

	return updatedNotification, nil
}
