package models

import (
	"time"

	"github.com/lib/pq"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Blog struct {
	ID        string         `gorm:"type:uuid;primary_key" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	AuthorID  string         `gorm:"type:uuid;not null" json:"author_id"`
	Author    User           `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Tags      pq.StringArray `gorm:"type:text[]" json:"tags,omitempty"`
	Images    pq.StringArray `gorm:"type:text[]" json:"images,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateBlogRequest struct {
	Title     string   `json:"title" validate:"required"`
	Content   string   `json:"content" validate:"required"`
	Tags      []string `json:"tags,omitempty"`
	ImageURLs []string `json:"image_urls,omitempty"`
}

func (b *Blog) Create(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &b)

	if err != nil {
		return err
	}

	return nil
}

func (b *Blog) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &b)

	if err != nil {
		return err
	}
	
	return nil
}

func (b *Blog) GetBlogById(db *gorm.DB, blogId string) (Blog, error) {
	var blog Blog
	err, nerr := postgresql.SelectOneFromDb(db, &blog, "id = ?", blogId)
	if nerr != nil {
		return blog, err
	}
	return blog, nil
}
