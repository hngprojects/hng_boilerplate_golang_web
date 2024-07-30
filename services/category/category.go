package category

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"gorm.io/gorm"
)

type Category struct {
	Name string `gorm:"size:255"`
}

func GetCategoryNames(db *gorm.DB, ctx *gin.Context) ([]Category, int, error) {
	var categories []Category

	ownerID, _ := middleware.GetIdFromToken(ctx)
	if ownerID == "" {
		return nil, http.StatusUnauthorized, errors.New("Unauthorized access")
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&Category{}).Select("name").Find(&categories).Error; err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, err
	}

	tx.Commit()

	return categories, http.StatusOK, nil
}
