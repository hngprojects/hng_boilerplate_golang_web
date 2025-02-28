package utility

import (
	"fmt"
	"strings"
	"gorm.io/gorm"
)

func IsUniqueSingleField[T any](db *gorm.DB, model T, field string, value interface{}) (bool, error) {
	var count int64

	err := db.Model(model).Where(fmt.Sprintf("%s = ?", field), value).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func IsUniqueMultipleFields(db *gorm.DB, model interface{}, fields map[string]interface{}) (bool, error) {
	var count int64

	queryParts := []string{}
	values := []interface{}{}

	// Iterate over the fields map to build query dynamically
	for field, value := range fields {
		queryParts = append(queryParts, fmt.Sprintf("%s = ?", field))
		values = append(values, value)
	}

	// Join all conditions using "AND"
	query := strings.Join(queryParts, " AND ")

	err := db.Model(model).Where(query, values...).Count(&count).Error
	fmt.Printf("Checking uniqueness on fields %+v with count %d\n", fields, count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
