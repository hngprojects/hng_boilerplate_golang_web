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
		Tags:     convertTags(req.Tags),
		Images:   convertImages(req.ImageURLs),
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

// convertTags is a helper function to convert a slice of tag names to Tag models
func convertTags(tagNames []string) []models.BlogTag {
	var tags []models.BlogTag
	for _, tagName := range tagNames {
		tags = append(tags, models.BlogTag{ID: utility.GenerateUUID(), Name: tagName})
	}
	return tags
}

// convertImages is a helper function to convert a slice of image URLs to Image models
func convertImages(imageURLs []string) []models.BlogImage {
	var images []models.BlogImage
	for _, imageURL := range imageURLs {
		images = append(images, models.BlogImage{ID: utility.GenerateUUID(), URL: imageURL})
	}
	return images
}
