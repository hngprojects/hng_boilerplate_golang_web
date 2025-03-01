package notifications

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/redis"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type NotificationObject struct {
	Notification *models.NotificationRecord
	ExtReq       request.ExternalRequest
	// rdb          *redis.Client
	rdb *redis.Redis
	Db  *postgresql.Postgresql
}

func NewNotificationObject(extReq request.ExternalRequest, rdb *redis.Redis, db *postgresql.Postgresql, notification *models.NotificationRecord) *NotificationObject {
	return &NotificationObject{
		ExtReq:       extReq,
		rdb:          rdb,
		Db:           db,
		Notification: notification,
	}
}

func ConvertToMapAndAddExtraData(data interface{}, newData map[string]interface{}) (map[string]interface{}, error) {
	var (
		mapData map[string]interface{}
	)

	mapData, err := utility.StructToMap(data)
	if err != nil {
		return mapData, err
	}

	for key, value := range newData {
		mapData[key] = value
	}

	return mapData, nil
}

func thisOrThatStr(this, that string) string {
	if this == "" {
		return that
	}
	return this
}
