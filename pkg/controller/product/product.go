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
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/product"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db             *storage.Database
	Validator      *validator.Validate
	Logger         *utility.Logger
	ExtReq         request.ExternalRequest
	ProductService product.ProductService
}

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

	owner_id, _ := middleware.GetIdFromToken(c)
	respData, code, err := base.ProductService.CreateProduct(req, owner_id)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("Product created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Product created successfully", respData)

	c.JSON(code, rd)
}

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
	ownerID, _ := middleware.GetIdFromToken(ctx)
	respData, code, err := base.ProductService.DeleteProduct(req, ownerID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		ctx.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("Product deleted successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Product deleted successfully", respData)

	ctx.JSON(code, rd)
}

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

	respData, code, err := base.ProductService.GetProduct(productId)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), "Product not found", nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("Product found successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Product found successfully", respData)

	c.JSON(code, rd)
}

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

	ownerID, _ := middleware.GetIdFromToken(c)
	respData, code, err := base.ProductService.UpdateProduct(req, ownerID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("Product updated successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Product updated successfully", respData)

	c.JSON(code, rd)
}

func (base *Controller) GetProductsInCategory(ctx *gin.Context) {
	category := ctx.Param("category")

	if category == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "", "Invalid category name", nil)
		ctx.JSON(http.StatusBadRequest, rd)
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	respData, code, err := base.ProductService.GetProductsInCategory(category, page, pageSize)
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
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	respData, code, err := base.ProductService.GetAllProducts(page, pageSize)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), "Products not found", nil)
		ctx.JSON(code, rd)
		return
	}

	base.Logger.Info("Products found successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Products found successfully", respData)

	ctx.JSON(code, rd)
}

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

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	respData, code, err := base.ProductService.FilterProducts(price, category, page, pageSize)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), "Products not found", nil)
		ctx.JSON(code, rd)
		return
	}

	base.Logger.Info("Products found successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Products found successfully", respData)

	ctx.JSON(code, rd)
}

func (base *Controller) UploadImage(ctx *gin.Context) {
	productId := ctx.Param("product_id")
	image, err := ctx.FormFile("image")

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		ctx.JSON(http.StatusBadRequest, rd)
		return
	}

	respData, code, err := base.ProductService.UploadImage(productId, image)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		ctx.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("Image uploaded successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Image uploaded successfully", respData)

	ctx.JSON(code, rd)
}
