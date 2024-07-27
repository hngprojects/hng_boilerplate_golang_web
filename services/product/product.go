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
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"owner_id":    product.OwnerID,
	}
	return responseData, http.StatusCreated, nil
}
func DeleteProduct(req string, db *gorm.DB, ctx *gin.Context) (gin.H, int, error) {
	var product models.Product
	if err := db.Where("id = ?", req).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, http.StatusNotFound, errors.New("product not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	ownerID, _ := middleware.GetIdFromToken(ctx)
	print(ownerID)

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
