package user

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func ReplaceUserRole(userID string, roleID int, db *gorm.DB) (gin.H, error) {

	var (
		user     = models.User{}
		role     = models.Role{}
		respData = gin.H{}
	)

	userExists := postgresql.CheckExists(db, &user, "id = ?", userID)
	if !userExists {
		return nil, errors.New("invalid user")
	}

	roleExists := postgresql.CheckExists(db, &role, "id = ?", roleID)
	if !roleExists {
		return nil, errors.New("invalid role")
	}

	userData, err := role.UpdateUserRole(db, userID, roleID)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	respData = gin.H{
		"username":   userData.Name,
		"first_name": userData.Profile.FirstName,
		"last_name":  userData.Profile.LastName,
		"phone":      userData.Profile.Phone,
		"role":       models.GetRoleName(models.RoleId(userData.Role)),
	}

	return respData, nil
}
