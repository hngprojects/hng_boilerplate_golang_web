package actions

import (
	"encoding/json"

	"github.com/go-redis/redis/v8"

	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/actions/names"
)

func AddNotificationToQueue(rdb *redis.Client, name names.NotificationName, data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	notificationRecord := models.NotificationRecord{
		Name: string(name),
		Data: string(dataByte),
		Sent: false,
	}

	// an instance of redis.Redis struct is passed to the PushToQueue method
	rd := inst.InitRed(rdb)
	err = notificationRecord.PushToQueue(rd)
	if err != nil {
		return err
	}

	return nil
}
