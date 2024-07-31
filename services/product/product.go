package product

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"

	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func CreateProduct(req models.CreateProductRequestModel, db *gorm.DB, c *gin.Context) (gin.H, int, error) {
	var (
		name         = strings.Title(strings.ToLower(req.Name))
		description  = req.Description
		price        = req.Price
		responseData gin.H
		categoryName = strings.Title(strings.ToLower(req.Category))
	)
	owner_id, _ := middleware.GetIdFromToken(c)

	tx := db.Begin()
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
		OwnerID:     owner_id,
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

func DeleteProduct(req models.DeleteProductRequestModel, db *gorm.DB, ctx *gin.Context) (gin.H, int, error) {
	var product models.Product
	if err := db.First(&product, req.ProductID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, http.StatusNotFound, errors.New("product not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	ownerID, err := middleware.GetIdFromToken(ctx)
	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("failed to get owner ID from token")
	}

	if product.OwnerID != ownerID {
		return nil, http.StatusForbidden, errors.New("you are not authorized to delete this product")
	}

	if err := db.Delete(&product).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	responseData := gin.H{
		"message": "Product deleted successfully",
	}
	return responseData, http.StatusOK, nil
}

func GetProduct(productId string, db *gorm.DB) (gin.H, int, error) {
	product := models.Product{}
	product, err := product.GetProduct(db, productId)
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

func UpdateProduct(req models.UpdateProductRequestModel, db *gorm.DB, ctx *gin.Context) (gin.H, int, error) {
	log.Printf("Received update request: %+v", req)
	var product models.Product
	if err := db.First(&product, "id = ?", req.ProductID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, http.StatusNotFound, errors.New("product not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	ownerID, _ := middleware.GetIdFromToken(ctx)

	if product.OwnerID != ownerID {
		return nil, http.StatusForbidden, errors.New("you are not authorized to update this product")
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price

	if err := db.Save(&product).Error; err != nil {
		log.Printf("Error saving product: %v", err) // Add this line
		return nil, http.StatusInternalServerError, err
	}

	responseData := gin.H{
		"message": "Product updated successfully",
	}
	return responseData, http.StatusOK, nil
}

func GetProductsInCategory(categoryName string, db *gorm.DB, c *gin.Context) (gin.H, int, error) {
	var category models.Category
	var products []models.Product

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	offset := (page - 1) * pageSize

	// Find category by name
	if err := db.Where("name = ?", categoryName).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, http.StatusNotFound, errors.New("category not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	// Get paginated products in category
	if err := db.Model(&category).Offset(offset).Limit(pageSize).Association("Products").Find(&products); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Debug information
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
