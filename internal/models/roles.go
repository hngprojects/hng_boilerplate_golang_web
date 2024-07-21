package models

type Role struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	RoleName       string       `gorm:"not null" json:"role_name"`
	OrganisationID uint         `gorm:"not null" json:"organisation_id"`
	Organisation   Organisation `gorm:"foreignKey:OrganisationID" json:"organisation"`
}

type Permission struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
}

type UserRole struct {
	UserID uint `gorm:"primaryKey" json:"user_id"`
	RoleID uint `gorm:"primaryKey" json:"role_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`
	Role   Role `gorm:"foreignKey:RoleID" json:"role"`
}

type RolePermission struct {
	RoleID       uint       `gorm:"primaryKey" json:"role_id"`
	PermissionID uint       `gorm:"primaryKey" json:"permission_id"`
	Role         Role       `gorm:"foreignKey:RoleID" json:"role"`
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"permission"`
}
