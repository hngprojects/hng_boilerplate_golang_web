package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"gorm.io/gorm"
	"net/http"
)

type RoleController struct {
	DB *gorm.DB
}

func NewRoleController(db *gorm.DB) *RoleController {
	return &RoleController{DB: db}
}

func (rc *RoleController) CreateRole(c *gin.Context) {
	var req models.CreateRoleRequestModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var org models.Organisation
	if err := rc.DB.First(&org, "id = ?", req.OrganizationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	permissions, err := models.GetPermissionsByIDs(rc.DB, req.PermissionIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		return
	}
	if len(permissions) != len(req.PermissionIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more permissions not found"})
		return
	}

	newUUID, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate UUID"})
		return
	}

	role := models.Role{
		ID:             newUUID.String(),
		RoleName:       req.RoleName,
		OrganizationID: req.OrganizationID,
	}

	if err := rc.DB.Transaction(func(tx *gorm.DB) error {
		if err := role.CreateRole(tx); err != nil {
			return err
		}
		if err := role.AddPermissionsToRole(tx, req.PermissionIDs); err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role and assign permissions"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Role created successfully",
		"status_code": http.StatusCreated,
	})
}
