package product

import (
	"net/http"
	"regexp"

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
