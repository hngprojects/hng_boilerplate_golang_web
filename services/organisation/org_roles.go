package organisation

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/user"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func CreateOrgRoles(req models.OrgRole, orgID string, db *gorm.DB, c *gin.Context) (gin.H, int, error) {
	var org models.Organisation

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	currentUserID, ok := userId.(string)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	currentUser, code, err := user.GetUser(currentUserID, db)
	if err != nil {
		return nil, code, err
	}

	orgData, err := org.CheckOrgExists(orgID, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gin.H{}, http.StatusNotFound, errors.New("organisation not found")
		}
		return gin.H{}, http.StatusBadRequest, err
	}

	isOwner, err := org.IsOwnerOfOrganisation(db, currentUser.ID, orgData.ID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if !isOwner {
		return nil, http.StatusForbidden, errors.New("not organization owner")
	}

	req.ID = utility.GenerateUUID()
	req.OrganisationID = orgData.ID

	if err := req.CreateOrgRole(db); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return gin.H{}, http.StatusConflict, errors.New("role name already exists")
		}
		return gin.H{}, http.StatusBadRequest, err
	}

	theResp := gin.H{
		"id":          req.ID,
		"name":        req.Name,
		"description": req.Description,
		"message":     "Role created successfully",
	}

	return theResp, http.StatusCreated, nil
}

func GetOrgRoles(db *gorm.DB, orgID string, c *gin.Context) ([]models.OrgRole, int, error) {
	var (
		org       models.Organisation
		role      models.OrgRole
		rolesData []models.OrgRole
	)

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	currentUserID, ok := userId.(string)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	currentUser, code, err := user.GetUser(currentUserID, db)
	if err != nil {
		return nil, code, err
	}

	orgData, err := org.CheckOrgExists(orgID, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New("organisation not found")
		}
		return nil, http.StatusBadRequest, err
	}

	isOwner, err := org.IsOwnerOfOrganisation(db, currentUser.ID, orgData.ID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if !isOwner {
		return nil, http.StatusForbidden, errors.New("not organization owner")
	}

	rolesData, err = role.GetOrgRoles(db, orgID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return rolesData, http.StatusOK, nil

}

func GetAOrgRole(db *gorm.DB, orgID, roleID string, c *gin.Context) (*models.OrgRole, int, error) {
	var (
		org       models.Organisation
		role      models.OrgRole
		rolesData models.OrgRole
	)

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	currentUserID, ok := userId.(string)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	currentUser, code, err := user.GetUser(currentUserID, db)
	if err != nil {
		return nil, code, err
	}

	orgData, err := org.CheckOrgExists(orgID, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New("organisation not found")
		}
		return nil, http.StatusBadRequest, err
	}

	isOwner, err := org.IsOwnerOfOrganisation(db, currentUser.ID, orgData.ID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if !isOwner {
		return nil, http.StatusForbidden, errors.New("not organization owner")
	}

	rolesData, err = role.GetAOrgRole(db, orgID, roleID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return &rolesData, http.StatusOK, nil

}

func DeleteOrgRole(db *gorm.DB, orgID, roleID string, c *gin.Context) (int, error) {
	var (
		org      models.Organisation
		role     models.OrgRole
		roleData models.OrgRole
	)

	userId, err := middleware.GetUserClaims(c, db, "user_id")
	if err != nil {
		return http.StatusNotFound, err
	}

	currentUserID, ok := userId.(string)
	if !ok {
		return http.StatusBadRequest, errors.New("user_id is not of type string")
	}

	currentUser, code, err := user.GetUser(currentUserID, db)
	if err != nil {
		return code, err
	}

	orgData, err := org.CheckOrgExists(orgID, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New("organisation not found")
		}
		return http.StatusBadRequest, err
	}

	isOwner, err := org.IsOwnerOfOrganisation(db, currentUser.ID, orgData.ID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if !isOwner {
		return http.StatusForbidden, errors.New("not organization owner")
	}

	roleData, err = role.GetAOrgRole(db, orgID, roleID)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = roleData.DeleteOrgRole(db)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil

}
