package services

// import (
// 	"errors"

// 	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
// 	"gorm.io/gorm"
// )

// type BlogService struct {
// 	Db *gorm.DB
// }

// func (s *BlogService) DeleteBlog(blogID string) error {
// 	// Assume Blog is your blog model
// 	var blog Blog
// 	if err := s.Db.First(&blog, "id = ?", blogID).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return errors.New("blog not found")
// 		}
// 		return err
// 	}

// 	return postgresql.DeleteRecordFromDb(s.Db, &blog)
// }
