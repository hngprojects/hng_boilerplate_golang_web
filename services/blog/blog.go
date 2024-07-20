package blogs

import (
	"errors"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func DeleteBlog(blogID string, db *gorm.DB) error {
	var blog models.Blog
	if err := db.First(&blog, "id = ?", blogID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("blog not found")
		}
		return err
	}

	return postgresql.DeleteRecordFromDb(db, &blog)
}
