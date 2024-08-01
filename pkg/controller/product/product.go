package product

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/product"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

// CreateProduct godoc
// @Summary Create Product
// @Description Create a new product
// @Tags Product
// @Accept json
// @Produce json
// @Param product body models.CreateProductRequestModel true "Product details"
// @Router /product [post]
func (base *Controller) CreateProduct(c *gin.Context) {

	var (
		req = models.CreateProductRequestModel{}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	respData, code, err := product.CreateProduct(req, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("Product created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Product created successfully", respData)

	c.JSON(code, rd)
}

// DeleteProductController godoc
// @Summary Delete Product
// @Description Delete a product
// @Tags Product
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} utility.Response
// @Failure 400,422 {object} utility.Response
// @Router /product/{product_id} [delete]
func (base *Controller) DeleteProductController(ctx *gin.Context) {
	var (
		req = models.DeleteProductRequestModel{}
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		ctx.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		ctx.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	respData, code, err := product.DeleteProduct(req, base.Db.Postgresql, ctx)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		ctx.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("Product deleted successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Product deleted successfully", respData)

	ctx.JSON(code, rd)
}

// GetProduct godoc
// @Summary Get Product
// @Description Get a Product
// @Tags Product
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} utility.Response
// @Failure 400,404,500 {object} utility.Response
// @Router /product/{product_id} [get]
func (base *Controller) GetProduct(c *gin.Context) {
	productId := c.Param("product_id")

	matched, err := regexp.MatchString("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", productId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), "An unexpected error occured", nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	if !matched {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "", "Invalid product ID", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	respData, code, err := product.GetProduct(productId, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), "Product not found", nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("Product found successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Product found successfully", respData)

	c.JSON(code, rd)
}

// UpdateProduct godoc
// @Summary Update Product
// @Description Update a product
// @Tags Product
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param product body models.UpdateProductRequestModel true "Product details"
// @Success 200 {object} utility.Response
// @Failure 400,422 {object} utility.Response
// @Router /product/{product_id} [put]
func (base *Controller) UpdateProduct(c *gin.Context) {
	var (
		req = models.UpdateProductRequestModel{}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	respData, code, err := product.UpdateProduct(req, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("Product updated successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Product updated successfully", respData)

	c.JSON(code, rd)
}

// GetProductsInCategory godoc
// @Summary Get Products In Category
// @Description Get all products in a category
// @Tags Product
// @Accept json
// @Produce json
// @Param category path string true "Category"
// @Success 200 {object} utility.Response
// @Failure 400,404 {object} utility.Response
// @Router /product/category/{category} [get]
func (base *Controller) GetProductsInCategory(ctx *gin.Context) {
	category := ctx.Param("category")

	if category == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "", "Invalid category name", nil)
		ctx.JSON(http.StatusBadRequest, rd)
		return
	}

	respData, code, err := product.GetProductsInCategory(category, base.Db.Postgresql, ctx)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), "Products not found", nil)
		ctx.JSON(code, rd)
		return
	}

	base.Logger.Info("Products found successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Products found successfully", respData)

	ctx.JSON(code, rd)
}

func (base *Controller) GetAllProducts(ctx *gin.Context) {
	respData, code, err := product.GetAllProducts(base.Db.Postgresql, ctx)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), "Products not found", nil)
		ctx.JSON(code, rd)
		return
	}

	base.Logger.Info("Products found successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Products found successfully", respData)

	ctx.JSON(code, rd)
}

// FilterProducts godoc
// @Summary Filter Products
// @Description Filter products by price and category
// @Tags Product
// @Accept json
// @Produce json
// @Param price query float64 true "Product Price"
// @Param category query string true "Product Category"
// @Success 200 {object} utility.Response
// @Failure 400,404 {object} utility.Response
// @Router /product/filter [get]
func (base *Controller) FilterProducts(ctx *gin.Context) {
	priceStr := ctx.Query("price")
	category := ctx.Query("category")

	if priceStr == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Price query parameter is required", "Invalid price", nil)
		ctx.JSON(http.StatusBadRequest, rd)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Println(err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), "Invalid price", nil)
		ctx.JSON(http.StatusBadRequest, rd)
		return
	}

	respData, code, err := product.FilterProducts(price, category, base.Db.Postgresql, ctx)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), "Products not found", nil)
		ctx.JSON(code, rd)
		return
	}

	base.Logger.Info("Products found successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Products found successfully", respData)

	ctx.JSON(code, rd)
}
