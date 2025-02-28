package models

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Blog struct {
	ID        string         `gorm:"type:uuid;primary_key" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	AuthorID  string         `gorm:"type:uuid;not null" json:"author_id"`
	Author    User           `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category  string         `gorm:"type:text" json:"category,omitempty"`
	Image     string         `gorm:"type:text" json:"image_url,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateBlogRequest struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Category string `json:"category,omitempty"`
	Image    string `json:"image_url,omitempty"`
}

type UpdateBlogRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Image    string `json:"image_url"`
}

func (b *Blog) Create(db database.DatabaseManager) error {
	err := db.CreateOneRecord(&b)

	if err != nil {
		return err
	}

	return nil
}

func (b *Blog) Delete(db database.DatabaseManager) error {
	err := db.DeleteRecordFromDb(&b)

	if err != nil {
		return err
	}

	return nil
}

func (b *Blog) GetBlogById(db database.DatabaseManager, blogId string) (Blog, error) {
	var blog Blog
	err, nerr := db.SelectOneFromDb(&blog, "id = ?", blogId)
	if nerr != nil {
		return blog, err
	}
	return blog, nil
}

func (b *Blog) GetAllBlogs(db database.DatabaseManager, c *gin.Context) ([]Blog, database.PaginationResponse, error) {
	var blog []Blog

	pagination := postgresql.GetPagination(c)

	paginationResponse, err := db.SelectAllFromDbOrderByPaginated(
		"created_at",
		"desc",
		"",
		pagination,
		&blog,
		nil,
	)

	if err != nil {
		return nil, paginationResponse, err
	}

	return blog, paginationResponse, nil
}

func (b *Blog) UpdateBlogById(db database.DatabaseManager, req UpdateBlogRequest, blogId string) (*Blog, error) {
	result, err := db.UpdateFields(&b, req, blogId)

	if err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("failed to update blog")
	}

	return b, nil
}

func (b *Blog) CheckBlogExists(blogId string, db database.DatabaseManager) (Blog, error) {
	blog, err := b.GetBlogById(db, blogId)
	if err != nil {
		return blog, err
	}

	return blog, nil
}
