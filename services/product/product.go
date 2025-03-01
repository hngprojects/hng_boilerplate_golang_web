package product

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type ProductService interface {
	CreateProduct(req models.CreateProductRequestModel, ownerID string) (gin.H, int, error)
	DeleteProduct(req models.DeleteProductRequestModel, ownerID string) (gin.H, int, error)
	GetProduct(productId string) (gin.H, int, error)
	UpdateProduct(req models.UpdateProductRequestModel, ownerID string) (gin.H, int, error)
	GetProductsInCategory(categoryName string, page, pageSize int) (gin.H, int, error)
	GetAllProducts(page, pageSize int) (gin.H, int, error)
	FilterProducts(price float64, category string, page, pageSize int) (gin.H, int, error)
	UploadImage(productID string, image *multipart.FileHeader) (gin.H, int, error)
}

type productService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) ProductService {
	return &productService{db: db}
}

func (s *productService) CreateProduct(req models.CreateProductRequestModel, ownerID string) (gin.H, int, error) {
	var (
		name         = strings.Title(strings.ToLower(req.Name))
		description  = req.Description
		price        = req.Price
		responseData gin.H
		categoryName = strings.Title(strings.ToLower(req.Category))
	)

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var category models.Category
	if err := tx.Where("LOWER(name) = LOWER(?)", categoryName).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If category does not exist, create it
			category = models.Category{
				ID:   utility.GenerateUUID(),
				Name: categoryName,
			}
			if err := tx.Create(&category).Error; err != nil {
				tx.Rollback()
				return nil, http.StatusInternalServerError, err
			}
		} else {
			tx.Rollback()
			return nil, http.StatusInternalServerError, err
		}
	}

	product := models.Product{
		ID:          utility.GenerateUUID(),
		Name:        name,
		Description: description,
		Price:       price,
		OwnerID:     ownerID,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, err
	}

	if err := tx.Model(&product).Association("Category").Append(&category); err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	responseData = gin.H{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"owner_id":    product.OwnerID,
		"category":    category.Name,
		"product_id":  product.ID,
	}
	return responseData, http.StatusCreated, nil
}

func (s *productService) DeleteProduct(req models.DeleteProductRequestModel, ownerID string) (gin.H, int, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var product models.Product
	if err := tx.Where("id = ?", req.ProductID).First(&product).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return nil, http.StatusNotFound, errors.New("product not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	if ownerID == "" {
		tx.Rollback()
		return nil, http.StatusUnauthorized, errors.New("failed to get owner ID from token")
	}

	if product.OwnerID != ownerID {
		tx.Rollback()
		return nil, http.StatusForbidden, errors.New("you are not authorized to delete this product")
	}

	if err := tx.Exec("DELETE FROM product_categories WHERE product_id = ?", product.ID).Error; err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, err
	}

	if err := tx.Delete(&product).Error; err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	responseData := gin.H{
		"message": "Product and its category associations deleted successfully",
	}
	return responseData, http.StatusOK, nil
}

func (s *productService) GetProduct(productId string) (gin.H, int, error) {
	product := models.Product{}
	product, err := product.GetProduct(s.db, productId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}

	responseData := gin.H{
		"id":          product.ID,
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"categories":  product.Category,
		"created_at":  product.CreatedAt,
		"updated_at":  product.UpdatedAt,
	}
	return responseData, http.StatusOK, nil
}

func (s *productService) UpdateProduct(req models.UpdateProductRequestModel, ownerID string) (gin.H, int, error) {
	log.Printf("Received update request: %+v", req)
	var product models.Product
	if err := s.db.First(&product, "id = ?", req.ProductID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, http.StatusNotFound, errors.New("product not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	if product.OwnerID != ownerID {
		return nil, http.StatusForbidden, errors.New("you are not authorized to update this product")
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price

	if err := s.db.Save(&product).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	responseData := gin.H{
		"message": "Product updated successfully",
	}
	return responseData, http.StatusOK, nil
}

func (s *productService) GetProductsInCategory(categoryName string, page, pageSize int) (gin.H, int, error) {
	var category models.Category
	var products []models.Product

	offset := (page - 1) * pageSize

	if err := s.db.Where("name = ?", categoryName).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, http.StatusNotFound, errors.New("category not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	if err := s.db.Model(&category).Offset(offset).Limit(pageSize).Association("Products").Find(&products); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	fmt.Printf("Category ID: %v\n", category.ID)
	fmt.Printf("Number of products found: %d\n", len(products))

	responseData := gin.H{
		"category":   categoryName,
		"products":   products,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int(math.Ceil(float64(len(category.Products)) / float64(pageSize))),
		"totalItems": len(category.Products),
	}
	return responseData, http.StatusOK, nil
}

func (s *productService) GetAllProducts(page, pageSize int) (gin.H, int, error) {
	var products []models.Product

	offset := (page - 1) * pageSize

	if err := s.db.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	responseData := gin.H{
		"products":   products,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int(math.Ceil(float64(len(products)) / float64(pageSize))),
		"totalItems": len(products),
	}
	return responseData, http.StatusOK, nil
}

func (s *productService) FilterProducts(price float64, category string, page, pageSize int) (gin.H, int, error) {
	var products []models.Product
	var totalCount int64

	query := s.db

	if price > 0 {
		query = query.Where("price <= ?", price)
	}

	if category != "" {
		query = query.Joins("JOIN product_categories ON products.id = product_categories.product_id").
			Joins("JOIN categories ON product_categories.category_id = categories.id").
			Where("categories.name = ?", category)
	}

	if err := query.Model(&models.Product{}).Count(&totalCount).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	offset := (page - 1) * pageSize

	if err := query.Order("price DESC").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	responseData := gin.H{
		"products":     products,
		"total_count":  totalCount,
		"current_page": page,
		"page_size":    pageSize,
		"total_pages":  int(math.Ceil(float64(totalCount) / float64(pageSize))),
	}
	return responseData, http.StatusOK, nil
}

func (s *productService) UploadImage(productID string, image *multipart.FileHeader) (gin.H, int, error) {
	product := models.Product{}
	if err := s.db.First(&product, "id = ?", productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gin.H{"error": "Product not found"}, http.StatusNotFound, err
		}
		return gin.H{"error": "Database error"}, http.StatusInternalServerError, err
	}

	if image == nil {
		return gin.H{"error": "No image file provided"}, http.StatusBadRequest, errors.New("no image file")
	}

	ext := filepath.Ext(image.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	if err := saveUploadedFile(image, fmt.Sprintf("images/%s", newFilename)); err != nil {
		return gin.H{"error": "Failed to save image"}, http.StatusInternalServerError, err
	}

	// Update product with new image filename
	product.Image = newFilename
	if err := s.db.Save(&product).Error; err != nil {
		return gin.H{"error": "Failed to update product"}, http.StatusInternalServerError, err
	}

	return gin.H{"message": "Image uploaded successfully"}, http.StatusOK, nil
}

// Helper function to save the uploaded file
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
