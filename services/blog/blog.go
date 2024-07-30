package service

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type BlogResponse struct {
	BlogID    string         `json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	ImageURLs pq.StringArray `json:"image_urls,omitempty"`
	Tags      pq.StringArray `json:"tags,omitempty"`
	Author    string         `json:"author"`
	CreatedAt time.Time      `json:"created_at"`
}

func CreateBlog(req models.CreateBlogRequest, db *gorm.DB, userId string) (BlogResponse, error) {
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
		return BlogResponse{}, err
	}

	user, err = user.GetUserByID(db, userId)

	if err != nil {
		return BlogResponse{}, err
	}

	response := BlogResponse{
		BlogID:    blog.ID,
		Title:     blog.Title,
		Content:   blog.Content,
		ImageURLs: blog.Images,
		Tags:      blog.Tags,
		Author:    user.Name,
		CreatedAt: blog.CreatedAt,
	}

	return response, nil
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

func GetBlogs(db *gorm.DB, c *gin.Context) ([]BlogResponse, postgresql.PaginationResponse, error) {
	var (
		blog models.Blog
		user models.User
	)
	blogs, paginationResponse, err := blog.GetAllBlogs(db, c)

	if err != nil {
		return nil, paginationResponse, err
	}

	var responses []BlogResponse

	for _, blog := range blogs {
		userId := blog.AuthorID
		user, _ = user.GetUserByID(db, userId)
		response := BlogResponse{
			BlogID:    blog.ID,
			Title:     blog.Title,
			Content:   blog.Content,
			ImageURLs: blog.Images,
			Tags:      blog.Tags,
			Author:    user.Name,
			CreatedAt: blog.CreatedAt,
		}

		responses = append(responses, response)
	}

	return responses, paginationResponse, nil
}

func GetBlogById(blogId string, db *gorm.DB) (BlogResponse, error) {
	var user models.User
	blog, err := CheckBlogExists(blogId, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return BlogResponse{}, errors.New("blog not found")
		}
		return BlogResponse{}, err
	}

	userId := blog.AuthorID
	user, _ = user.GetUserByID(db, userId)

	response := BlogResponse{
		BlogID:    blog.ID,
		Title:     blog.Title,
		Content:   blog.Content,
		ImageURLs: blog.Images,
		Tags:      blog.Tags,
		Author:    user.Name,
		CreatedAt: blog.CreatedAt,
	}

	return response, nil
}

func CheckBlogExists(blogId string, db *gorm.DB) (models.Blog, error) {
	var blog models.Blog

	blog, err := blog.GetBlogById(db, blogId)
	if err != nil {
		return blog, err
	}

	return blog, nil
}
