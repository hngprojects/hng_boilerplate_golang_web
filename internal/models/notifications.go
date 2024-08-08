package models

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
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

func (n *Notification) CreateNotification(db *gorm.DB) (Notification, error) {

	err := postgresql.CreateOneRecord(db, &n)

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

func (n *Notification) GetNotificationByID(db *gorm.DB, ID string) (Notification, error) {
	var notification Notification

	err, _ := postgresql.SelectOneFromDb(db, &notification, "id = ?", ID)
	if err != nil {
		return notification, err
	}
	return notification, nil
}

func (n *Notification) FetchAllNotifications(db *gorm.DB, c *gin.Context) ([]Notification, map[string]int64, error) {
	var notifications []Notification
	type additionalData map[string]int64

	err := postgresql.SelectAllFromDb(db, "", &notifications, "")
	if err != nil {
		return nil, additionalData{}, err
	}

	totalCount, err := postgresql.CountRecords(db, &notifications)
	if err != nil {
		return nil, additionalData{}, err
	}

	unreadCount, err := postgresql.CountSpecificRecords(db, &notifications, "is_read = false")
	if err != nil {
		return nil, additionalData{}, err
	}

	data := additionalData{
		"total_count":  totalCount,
		"unread_count": unreadCount,
	}

	return notifications, data, nil
}

func (n *Notification) FetchUnReadNotifications(db *gorm.DB, c *gin.Context) ([]Notification, map[string]int64, error) {
	var notifications []Notification
	type additionalData map[string]int64

	totalCount, err := postgresql.CountRecords(db, &n)
	if err != nil {
		return nil, additionalData{}, err
	}

	unreadCount, err := postgresql.CountSpecificRecords(db, &n, "is_read = false")
	if err != nil {
		return nil, additionalData{}, err
	}

	data := additionalData{
		"total_count":  totalCount,
		"unread_count": unreadCount,
	}

	err = postgresql.SelectAllFromDb(db, "", &notifications, "is_read = ?", false)
	if err != nil {
		return nil, additionalData{}, err
	}
	return notifications, data, nil
}

func (n *Notification) UpdateNotification(db *gorm.DB, notifReq UpdateNotificationReq, ID string) (*Notification, error) {

	exists := postgresql.CheckExists(db, &n, "id = ?", ID)
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	oldNotification, err := n.GetNotificationByID(db, ID)
	if err != nil {
		return nil, err
	}

	res, err := postgresql.UpdateFields(db, &n, notifReq, ID)
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

func (n *Notification) DeleteNotificationByUserID(db *gorm.DB, ID string) error {
	var notifications []Notification

	exists := postgresql.CheckExists(db, &notifications, "user_id = ?", ID)
	if !exists {
		return gorm.ErrRecordNotFound
	}

	err := postgresql.SelectAllFromDb(db, "", &notifications, "user_id = ?", ID)
	if err != nil {
		return err
	}

	err = postgresql.DeleteRecordFromDb(db, &notifications)
	if err != nil {
		return err
	}
	return nil
}

func (n *NotificationSettings) GetNotificationSettingsByID(db *gorm.DB, ID string) (NotificationSettings, error) {
	var notificationSettings NotificationSettings

	err, nerr := postgresql.SelectOneFromDb(db, &notificationSettings, "user_id = ?", ID)
	if err != nil {
		return notificationSettings, nerr
	}
	return notificationSettings, nil
}

func (n *NotificationSettings) CreateNotificationSettings(db *gorm.DB) (NotificationSettings, error) {
	err := postgresql.CreateOneRecord(db, &n)

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

func (n *NotificationSettings) UpdateNotificationSettings(db *gorm.DB, ID string) (NotificationSettings, error) {
	n.ID = ID

	exists := postgresql.CheckExists(db, &NotificationSettings{}, "user_id = ?", ID)
	if !exists {
		return NotificationSettings{}, gorm.ErrRecordNotFound
	}

	_, err := postgresql.SaveAllFields(db, n)
	if err != nil {
		return NotificationSettings{}, err
	}

	updatedNotification := NotificationSettings{}
	err = db.First(&updatedNotification, "id = ?", ID).Error
	if err != nil {
		return NotificationSettings{}, err
	}

	return updatedNotification, nil
}
