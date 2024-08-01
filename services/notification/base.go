package notifications

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)


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