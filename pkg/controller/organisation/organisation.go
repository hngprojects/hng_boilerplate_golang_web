package organisation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	service "github.com/hngprojects/hng_boilerplate_golang_web/services/organisation"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) CreateOrganisation(c *gin.Context) {

	var (
		req = models.CreateOrgRequestModel{}
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

	reqData, code, err := service.ValidateCreateOrgRequest(req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)

	userId := userClaims["user_id"].(string)

	respData, err := service.CreateOrganisation(reqData, base.Db.Postgresql, userId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("organisation created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "organisation created successfully", respData)

	c.JSON(http.StatusCreated, rd)
}

func (base *Controller) GetOrganisationById(c *gin.Context) {
	orgId := c.Param("org_id")

	if _, err := uuid.Parse(orgId); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid organisation id format", "failed to delete organisation", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", "failed to delete organisation", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	orgData, err := service.GetOrganisationById(orgId, userId, base.Db.Postgresql)

	if err != nil {
		switch err.Error() {
		case "organisation not found":
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), "failed to retrieve organisation", nil)
			c.JSON(http.StatusNotFound, rd)
		case "user not authorised to retrieve this organisation":
			rd := utility.BuildErrorResponse(http.StatusForbidden, "error", err.Error(), "failed to retrieve organisation", nil)
			c.JSON(http.StatusForbidden, rd)
		default:
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "failed to retrieve organisation", err.Error(), nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("organisation retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "organisation retrieved successfully", orgData)

	c.JSON(http.StatusOK, rd)

}

func (base *Controller) UpdateOrganisation(c *gin.Context) {
	orgId := c.Param("org_id")
	var updateReq models.CreateOrgRequestModel

	if _, err := uuid.Parse(orgId); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid organisation id format", "failed to update organisation", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", "failed to update organisation", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	if err := c.ShouldBind(&updateReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&updateReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	updatedOrg, err := service.UpdateOrganisation(orgId, userId, updateReq, base.Db.Postgresql)

	if err != nil {
		switch err.Error() {
		case "organisation not found":
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), "failed to update organisation", nil)
			c.JSON(http.StatusNotFound, rd)
		case "user not authorised to update this organisation":
			rd := utility.BuildErrorResponse(http.StatusForbidden, "error", err.Error(), "failed to update organisation", nil)
			c.JSON(http.StatusForbidden, rd)
		case "organisation already exists with the given email":
			rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), "failed to update organisation", nil)
			c.JSON(http.StatusForbidden, rd)
		default:
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "failed to update organisation", err.Error(), nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("organisation updated successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "organisation updated successfully", updatedOrg)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) DeleteOrganisation(c *gin.Context) {
	orgId := c.Param("org_id")

	if _, err := uuid.Parse(orgId); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid organisation id format", "failed to delete organisation", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", "failed to delete organisation", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	if err := service.DeleteOrganisation(orgId, userId, base.Db.Postgresql); err != nil {
		switch err.Error() {
		case "organisation not found":
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), "failed to delete organisation", nil)
			c.JSON(http.StatusNotFound, rd)
		case "user not authorised to delete this organisation":
			rd := utility.BuildErrorResponse(http.StatusForbidden, "error", err.Error(), "failed to delete organisation", nil)
			c.JSON(http.StatusForbidden, rd)
		default:
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "failed to delete organisation", err.Error(), nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("organisation deleted successfully")
	rd := utility.BuildSuccessResponse(http.StatusNoContent, "", nil)
	c.JSON(http.StatusNoContent, rd)
}
