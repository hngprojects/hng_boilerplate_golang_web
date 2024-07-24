package product

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"

	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func CreateProduct(req models.CreateProductRequestModel, db *gorm.DB) (gin.H, int, error) {
	var (
		name = strings.Title(strings.ToLower(req.Name))
		description = req.Description
		price = req.Price
		owner_id = req.OwnerID
		responseData gin.H
	)

	product := models.Product{
		ID:       		utility.GenerateUUID(),
		Name:     		name,
		Description:    description,
		Price: 			price,
		OwnerID: 		owner_id,
	}

	err := product.CreateProduct(db)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	responseData = gin.H{
		"name":        product.Name,
		"description": product.Description,
		"price":   	   product.Price,
		"owner_id":    product.OwnerID,
	}
	return responseData, http.StatusCreated, nil
}