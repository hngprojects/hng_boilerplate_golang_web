package service

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

// BlogService defines the interface for blog-related operations

type BlogService interface {
	CreateBlog(req models.CreateBlogRequest, userId string) (BlogResponse, error)
	DeleteBlog(blogId string, userId string) error
	GetBlogs(c *gin.Context) ([]BlogResponse, database.PaginationResponse, error)
	GetBlogById(blogId string) (BlogResponse, error)
	UpdateBlogById(blogId string, userId string, req models.UpdateBlogRequest) (BlogResponse, error)
}

// BlogServiceImpl is the concrete implementation of BlogService
type BlogServiceImpl struct {
	db database.DatabaseManager
}

// NewBlogService creates a new BlogService instance
func NewBlogService(db database.DatabaseManager) BlogService {
	return &BlogServiceImpl{db: db}
}

type BlogResponse struct {
	BlogID    string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Image     string    `json:"image_url,omitempty"`
	Category  string    `json:"category,omitempty"`
	Author    string    `json:"author"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *BlogServiceImpl) CreateBlog(req models.CreateBlogRequest, userId string) (BlogResponse, error) {
	// instance of Postgresql db
	pdb := s.db // No need for `inst.InitDB(db)` anymore

	var user models.User
	blog := models.Blog{
		ID:       utility.GenerateUUID(),
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: userId,
		Category: req.Category,
		Image:    req.Image,
	}

	err := blog.Create(pdb)

	if err != nil {
		return BlogResponse{}, err
	}

	user, err = user.GetUserByID(pdb, userId)

	if err != nil {
		return BlogResponse{}, err
	}

	response := BlogResponse{
		BlogID:    blog.ID,
		Title:     blog.Title,
		Content:   blog.Content,
		Image:     blog.Image,
		Category:  blog.Category,
		Author:    user.Name,
		AuthorID:  user.ID,
		CreatedAt: blog.CreatedAt,
	}

	return response, nil
}

func (s *BlogServiceImpl) DeleteBlog(blogId string, userId string) error {
	// instance of Postgresql db
	pdb := s.db
	var blog models.Blog
	blog, err := blog.CheckBlogExists(blogId, pdb)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("blog not found")
		}
		return err
	}

	if blog.AuthorID != userId {
		return errors.New("user not authorised to delete blog")
	}

	return blog.Delete(pdb)
}

func (s *BlogServiceImpl) GetBlogs(c *gin.Context) ([]BlogResponse, database.PaginationResponse, error) {
	// instance of Postgresql db
	pdb := s.db
	var (
		blog models.Blog
		user models.User
	)
	blogs, paginationResponse, err := blog.GetAllBlogs(pdb, c)

	if err != nil {
		return nil, paginationResponse, err
	}

	var responses []BlogResponse

	for _, blog := range blogs {
		userId := blog.AuthorID
		user, _ = user.GetUserByID(pdb, userId)
		response := BlogResponse{
			BlogID:    blog.ID,
			Title:     blog.Title,
			Content:   blog.Content,
			Image:     blog.Image,
			Category:  blog.Category,
			Author:    user.Name,
			AuthorID:  user.ID,
			CreatedAt: blog.CreatedAt,
		}

		responses = append(responses, response)
	}

	return responses, paginationResponse, nil
}

func (s *BlogServiceImpl) GetBlogById(blogId string) (BlogResponse, error) {
	// instance of Postgresql db
	pdb := s.db
	var (
		user models.User
		blog models.Blog
	)
	blog, err := blog.CheckBlogExists(blogId, pdb)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return BlogResponse{}, errors.New("blog not found")
		}
		return BlogResponse{}, err
	}

	userId := blog.AuthorID
	user, _ = user.GetUserByID(pdb, userId)

	response := BlogResponse{
		BlogID:    blog.ID,
		Title:     blog.Title,
		Content:   blog.Content,
		Image:     blog.Image,
		Category:  blog.Category,
		Author:    user.Name,
		AuthorID:  user.ID,
		CreatedAt: blog.CreatedAt,
	}

	return response, nil
}

func (s *BlogServiceImpl) UpdateBlogById(blogId string, userId string, req models.UpdateBlogRequest) (BlogResponse, error) {
	// instance of Postgresql db
	pdb := s.db
	var (
		user models.User
		blog models.Blog
	)
	blog, err := blog.CheckBlogExists(blogId, pdb)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return BlogResponse{}, errors.New("blog not found")
		}
		return BlogResponse{}, err
	}

	if blog.AuthorID != userId {
		return BlogResponse{}, errors.New("user not authorised to update blog")
	}

	user, _ = user.GetUserByID(pdb, userId)

	updatedBlog, err := blog.UpdateBlogById(pdb, req, blogId)

	if err != nil {
		return BlogResponse{}, err
	}

	response := BlogResponse{
		BlogID:    updatedBlog.ID,
		Title:     updatedBlog.Title,
		Content:   updatedBlog.Content,
		Image:     updatedBlog.Image,
		Category:  updatedBlog.Category,
		Author:    user.Name,
		AuthorID:  userId,
		UpdatedAt: updatedBlog.UpdatedAt,
	}

	return response, nil
}
