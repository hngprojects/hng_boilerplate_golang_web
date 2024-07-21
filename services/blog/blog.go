package service

import (
	"errors"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func CreateBlog(req models.CreateBlogRequest, db *gorm.DB, userId string) (*models.Blog, error) {
	var existingBlog models.Blog
	blog := models.Blog{
		ID:       utility.GenerateUUID(),
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: userId,
		Tags:     req.Tags,
		Images:   req.ImageURLs,
	}

	if err := db.Where("title = ?", &blog.Title).First(&existingBlog).Error; err == nil {
		return nil, errors.New("blog with this title already exists")
	}

	err := blog.Create(db)

	if err != nil {
		return nil, err
	}

	return &blog, nil

}

func DeleteBlog(blogID string, db *gorm.DB) error {
	var blog models.Blog
	if err := db.First(&blog, "id = ?", blogID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("blog not found")
		}
		return err
	}

	return blog.Delete(db)

}
