package cronjobs

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/actions"
)

func SendNotifications(extReq request.ExternalRequest, db storage.Database) {
	notificationRecord := models.NotificationRecord{}

	res, err := notificationRecord.PopFromQueue(db.Redis)

	if err != nil {
		extReq.Logger.Error("error getting notificatin records: ", err.Error())
		return
	}

	extReq.Logger.Error("Sending records found: ", res)

	err = actions.Send(extReq, db.Postgresql, db.Redis, &res)

	if err != nil {
		extReq.Logger.Error("error getting notificatin records: ", err.Error())
		return
	}

}
