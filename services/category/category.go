package category

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"gorm.io/gorm"
)

type Category struct {
	Name string `gorm:"size:255"`
}

type PaginatedResponse struct {
	Categories []Category `json:"categories"`
	TotalCount int64      `json:"totalCount"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
}

// CategoryService defines the interface for category operations
type CategoryService interface {
	GetCategoryNames(ctx *gin.Context) (*PaginatedResponse, int, error)
}

// categoryService implements CategoryService
type categoryService struct {
	db *gorm.DB
}

// NewCategoryService initializes a new CategoryService instance
func NewCategoryService(db *gorm.DB) CategoryService {
	return &categoryService{db: db}
}

// GetCategoryNames fetches paginated category names
func (s *categoryService) GetCategoryNames(ctx *gin.Context) (*PaginatedResponse, int, error) {
	ownerID, _ := middleware.GetIdFromToken(ctx)
	if ownerID == "" {
		return nil, http.StatusUnauthorized, errors.New("unauthorized access")
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var categories []Category
	var totalCount int64

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&Category{}).Count(&totalCount).Error; err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, err
	}

	if err := tx.Model(&Category{}).
		Select("name").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&categories).Error; err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, err
	}

	tx.Commit()

	response := &PaginatedResponse{
		Categories: categories,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	return response, http.StatusOK, nil
}
