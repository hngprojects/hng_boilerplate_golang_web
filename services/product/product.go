package product

import (
	"errors"
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
		"id":          product.ID,
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"owner_id":    product.OwnerID,
	}
	return responseData, http.StatusCreated, nil
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
	return responseData, http.StatusCreated, nil
}
