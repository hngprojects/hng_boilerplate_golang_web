package models

import (
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	MemberID string         `gorm:"unique" json:"memberId"`
	TeamID   string         `json:"teamId" binding:"required"`
	Role     string         `json:"role" binding:"required"`
	Teams    []Organisation `gorm:"many2many:team_members"`
}
