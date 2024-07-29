package product

import (
	"errors"
	"log"
	"net/http"
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
	)

	owner_id, _ := middleware.GetIdFromToken(c)

	product := models.Product{
		ID:          utility.GenerateUUID(),
		Name:        name,
		Description: description,
		Price:       price,
		OwnerID:     owner_id,
	}

	err := product.CreateProduct(db)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	responseData = gin.H{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"owner_id":    product.OwnerID,
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
