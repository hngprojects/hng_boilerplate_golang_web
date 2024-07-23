package service

import (
	"errors"
	"fmt"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func ReplaceUserRole(userID string, roleID int, db *gorm.DB) (*models.User, error) {

	user := models.User{}
	role := models.Role{}

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

	return userData, nil
}
