package role

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type RoleController struct {
	Db        *gorm.DB
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

type CreateRoleRequest struct {
	RoleName       string   `json:"role_name" binding:"required"`
	OrganizationID string   `json:"organization_id" binding:"required"`
	PermissionIDs  []string `json:"permission_ids" binding:"required"`
}

func (controller *RoleController) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "status_code": 400})
		return
	}

	// Authorization check (assuming middleware sets user role)
	userRole := c.GetString("userRole")
	if userRole != "admin" && userRole != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden", "status_code": 403})
		return
	}

	// Role creation logic
	role := models.Role{
		ID:             utility.GenerateUUID(), // Assuming you have a function to generate UUIDs
		RoleName:       req.RoleName,
		OrganizationID: req.OrganizationID,
	}

	if err := controller.Db.Create(&role).Error; err != nil {
		controller.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Role creation failed", "status_code": 500})
		return
	}

	// Adding permissions
	for _, permID := range req.PermissionIDs {
		rolePermission := models.RolePermission{
			RoleID:       role.ID,
			PermissionID: permID,
		}
		controller.Db.Create(&rolePermission)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Role created successfully", "status_code": 201})
}
