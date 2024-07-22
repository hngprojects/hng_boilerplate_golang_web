package member

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type Controller struct {
	Db        *storage.Database
	Logger    *utility.Logger
}

type ChangeRoleRequest struct {
	NewRole string `json:"newRole" binding:"required,oneof=admin user superadmin member"`
}

func (base *Controller) ChangeUserRole(c *gin.Context) {
	var request ChangeRoleRequest
	var member models.Member
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamID := c.Param("teamId")
	memberID := c.Param("memberId")

	if err := base.Db.Postgresql.Where("team_id = ? AND member_id = ?", teamID, memberID).First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	member.Role = request.NewRole
	if err := base.Db.Postgresql.Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	base.Logger.Info("organisation created successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":  "Team member role updated successfully",
		"teamId":   member.TeamID,
		"memberId": member.MemberID,
		"newRole":  request.NewRole,
	})
}
