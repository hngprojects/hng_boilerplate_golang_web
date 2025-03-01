package user

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func (s *userService) ReplaceUserRole(userID string, roleID int) (gin.H, error) {

	var (
		user     = models.User{}
		role     = models.Role{}
		respData = gin.H{}
	)
	pdb := inst.InitDB(s.db)

	userExists := pdb.CheckExists(&user, "id = ?", userID)
	if !userExists {
		return nil, errors.New("invalid user")
	}

	roleExists := pdb.CheckExists(&role, "id = ?", roleID)
	if !roleExists {
		return nil, errors.New("invalid role")
	}

	userData, err := role.UpdateUserRole(pdb, userID, roleID)
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
