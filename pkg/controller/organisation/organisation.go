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

func (base *Controller) GetOrganisation(c *gin.Context) {
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

	orgData, err := service.GetOrganisation(orgId, userId, base.Db.Postgresql)

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
	var updateReq models.UpdateOrgRequestModel

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

func (base *Controller) AddUserToOrganisation(c *gin.Context) {
	orgId := c.Param("org_id")

	var req models.AddUserToOrgRequestModel
	if err := c.ShouldBind(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if _, err := uuid.Parse(orgId); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid organisation id format", "failed to update organisation", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err := base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	err = service.AddUserToOrganisation(orgId, req, base.Db.Postgresql)

	if err != nil {
		switch err.Error() {
		case "organisation not found":
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), "failed to add user to organisation", nil)
			c.JSON(http.StatusNotFound, rd)
		case "user not found":
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), "failed to add user to organisation", nil)
			c.JSON(http.StatusNotFound, rd)
		case "user already added to organisation":
			rd := utility.BuildErrorResponse(http.StatusConflict, "error", err.Error(), "failed to add user to organisation", nil)
			c.JSON(http.StatusNotFound, rd)
		default:
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "failed to add user to organisation", err.Error(), nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("user added to organisation successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "user added to organisation successfully", nil)

	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetUsersInOrganisation(c *gin.Context) {
	orgId := c.Param("org_id")

	if _, err := uuid.Parse(orgId); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid organisation id format", "failed to retrieve users", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", "failed to retrieve users", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	users, paginationResponse, err := service.GetUsersInOrganisation(orgId, userId, base.Db.Postgresql, c)

	if err != nil {
		switch err.Error() {
		case "organisation not found":
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), "failed to retrieve users", nil)
			c.JSON(http.StatusNotFound, rd)
		case "user does not have access to the organisation":
			rd := utility.BuildErrorResponse(http.StatusForbidden, "error", err.Error(), "failed to retrieve users", nil)
			c.JSON(http.StatusNotFound, rd)
		default:
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "failed to retrieve users", err.Error(), nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	paginationData := map[string]interface{}{
		"current_page": paginationResponse.CurrentPage,
		"total_pages":  paginationResponse.TotalPagesCount,
		"page_size":    paginationResponse.PageCount,
		"total_items":  len(users),
	}

	base.Logger.Info("users retrieved successfully")
	response := utility.BuildSuccessResponse(http.StatusOK, "users retrieved successfully", users, paginationData)

	c.JSON(http.StatusOK, response)
}
