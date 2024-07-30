package service

import (
	"errors"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"github.com/lib/pq"
	"gorm.io/gorm"
)


func CreateBlog(req models.CreateBlogRequest, db *gorm.DB, userId string) error {
	var user models.User
	blog := models.Blog{
		ID:       utility.GenerateUUID(),
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: userId,
		Tags:     pq.StringArray(req.Tags),
		Images:   pq.StringArray(req.ImageURLs),
	}

	err := blog.Create(db)

	if err != nil {
		return err
	}

	user, err = user.GetUserByID(db, userId)

	if err != nil {
		return err
	}

	return nil

}

func DeleteBlog(blogId string, db *gorm.DB) error {
	blog, err := CheckBlogExists(blogId, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("blog not found")
		}
		return err
	}

	return blog.Delete(db)
}

func CheckBlogExists(blogId string, db *gorm.DB) (models.Blog, error) {
	var blog models.Blog

	blog, err := blog.GetBlogById(db, blogId)
	if err != nil {
		return blog, err
	}

	return blog, nil
}
