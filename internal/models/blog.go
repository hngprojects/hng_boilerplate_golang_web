package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Blog struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	AuthorID  string    `gorm:"type:uuid;not null" json:"author_id"`
	Author    User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Tags      []string  `gorm:"type:text[]" json:"tags,omitempty"`
	Images    []string  `gorm:"type:text[]" json:"images,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

type CreateBlogRequest struct {
	Title     string   `json:"title" validate:"required"`
	Content   string   `json:"content" validate:"required"`
	Tags      []string `json:"tags,omitempty"`
	ImageURLs []string `json:"image_urls,omitempty"`
}

func (blog *Blog) Create(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, blog)

	if err != nil {
		return err
	}

	return nil
}

func (blog *Blog) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, blog)

	if err != nil {
		return err
	}

	return nil
}
