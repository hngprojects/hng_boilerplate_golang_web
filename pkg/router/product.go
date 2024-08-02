package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/product"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Product(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	product := product.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	productUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql))
	{
		productUrl.POST("/products", product.CreateProduct)
		productUrl.DELETE("/products", product.DeleteProductController)
		productUrl.GET("/products/:product_id", product.GetProduct)
		productUrl.PUT("/products/", product.UpdateProduct)
		productUrl.GET("/products/categories/:category", product.GetProductsInCategory)
		productUrl.GET("/products", product.GetAllProducts)
		productUrl.GET("/products/filter/", product.FilterProducts)
		productUrl.PATCH("/products/image/:product_id", product.UploadImage)
	}

	return r
}
