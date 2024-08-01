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

func UpdateOrgRoles(req models.OrgRole, orgID, roleID string, db *gorm.DB, c *gin.Context) (gin.H, int, error) {
	var (
		org      models.Organisation
		roleData models.OrgRole
		role     models.OrgRole
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

	roleData, err = role.GetAOrgRole(db, orgID, roleID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	roleData.Name = req.Name
	roleData.Description = req.Description

	if err := roleData.UpdateOrgRole(db); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return gin.H{}, http.StatusConflict, errors.New("role name already exists")
		}
		return gin.H{}, http.StatusBadRequest, err
	}

	theResp := gin.H{
		"id":          roleData.ID,
		"name":        roleData.Name,
		"description": roleData.Description,
	}

	return theResp, http.StatusOK, nil
}

func UpdateOrgPermissions(req models.Permission, orgID, roleID string, db *gorm.DB, c *gin.Context) (int, error) {
	var (
		org      models.Organisation
		roleData models.OrgRole
		role     models.OrgRole
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

	req.ID = utility.GenerateUUID()
	req.RoleID = roleData.ID

	if err := req.UpdateOrgPermissions(db); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return http.StatusConflict, errors.New("permission already exists")
		}
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
