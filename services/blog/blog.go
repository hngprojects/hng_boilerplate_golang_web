package service

import (
	"errors"
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CreateResponse struct {
	BlogID    string         `json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	ImageURLs pq.StringArray `json:"image_urls"`
	Tags      pq.StringArray `json:"tags"`
	Author    string         `json:"author"`
	CreatedAt time.Time      `json:"created_at"`
}

func CreateBlog(req models.CreateBlogRequest, db *gorm.DB, userId string) (CreateResponse, error) {
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
		return CreateResponse{}, err
	}

	user, err = user.GetUserByID(db, userId)

	if err != nil {
		return CreateResponse{}, err
	}

	createResponse := CreateResponse{
		BlogID:    blog.ID,
		Title:     blog.Title,
		Content:   blog.Content,
		ImageURLs: blog.Images,
		Tags:      blog.Tags,
		Author:    user.Name,
		CreatedAt: blog.CreatedAt,
	}

	return createResponse, nil
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
