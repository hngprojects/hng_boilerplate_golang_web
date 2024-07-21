package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Blog struct {
	ID        string      `gorm:"type:uuid;primary_key" json:"id"`
	Title     string      `gorm:"not null" json:"title"`
	Content   string      `gorm:"type:text" json:"content"`
	AuthorID  string      `gorm:"type:uuid;not null" json:"author_id"`
	Author    User        `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Tags      []BlogTag   `gorm:"many2many:blog_tags" json:"tags,omitempty"`
	Images    []BlogImage `gorm:"foreignKey:BlogID" json:"images,omitempty"`
	CreatedAt time.Time   `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time   `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

type BlogTag struct {
	ID    string `gorm:"type:uuid;primary_key" json:"id"`
	Name  string `gorm:"not null;unique" json:"name"`
	Blogs []Blog `gorm:"many2many:blog_tags" json:"blogs,omitempty"`
}

type BlogImage struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	URL       string    `gorm:"not null" json:"url"`
	BlogID    string    `gorm:"type:uuid;not null" json:"blog_id"`
	CreatedAt time.Time `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
}

type CreateBlogRequest struct {
	Title     string   `json:"title" validate:"required"`
	Content   string   `json:"content" validate:"required"`
	Tags      []string `json:"tags,omitempty"`
	ImageURLs []string `json:"image_urls,omitempty"`
}

func (blog *Blog) Create(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &blog)

	if err != nil {
		return err
	}

	return nil
}

func (blog *Blog) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &blog)

	if err != nil {
		return err
	}

	return nil
}
