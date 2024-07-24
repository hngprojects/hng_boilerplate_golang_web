package invite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}
type InvitationRequest struct {
	Emails []string `json:"emails" validate:"required"`
	OrgID  string   `json:"org_id" binding:"uuid"`
}



// Invitation model definition


func (base *Controller) PostInvite(c *gin.Context) {

	var inviteReq InvitationRequest

	if err := c.ShouldBindJSON(&inviteReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err := base.Validator.Struct(&inviteReq)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest,"error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// if err != nil {
	// 	rd := utility.BuildErrorResponse(http.StatusNotFound, "error", err.Error(), err, nil)
	// 	c.JSON(http.StatusInternalServerError, rd)
	// 	return
	// }

	base.Logger.Info("invite posted successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "", "invite url")

	c.JSON(http.StatusOK, rd)
}


// InvLink struct for binding the request body
type InvLink struct {
	InvitationLink string `json:"invitation_link"`
}

// DeactivateInvitation handler
func (base *Controller) DeactivateInvitation (ctx *gin.Context) {

	// Bind the request body to the invLink struct
	var invLink InvLink
	if err := ctx.Bind(&invLink); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"errors": []gin.H{
				{
					"field":   "invitation_link",
					"message": "Invalid request body",
				},
			},
			"status_code": 400,
		})
		return
	}

	db := storage.Connection()
	var invitation models.Invitation 
	result := db.Postgresql.Where("invitation_link = ?", invLink.InvitationLink).First(&invitation)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Invitation link not found",
				"status_code": 404,
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database error",
				"errors": []gin.H{
					{
						"field":   "database",
						"message": result.Error.Error(),
					},
				},
				"status_code": 500,
			})
		}
		return
	}

	// Update the isValid field to false
	invitation.IsValid = false
	if err := db.Postgresql.Save(&invitation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to deactivate invitation link",
			"errors": []gin.H{
				{
					"field":   "database",
					"message": err.Error(),
				},
			},
			"status_code": 500,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Invitation link deactivated successfully",
		"status_code": 200,
	}) 
}
