package models

type Role struct {
	ID             string       `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	RoleName       string       `gorm:"type:varchar(255);not null" json:"role_name"`
	OrganizationID string       `gorm:"type:uuid;not null"`
	Permissions    []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	ID   string `gorm:"type:uuid;primaryKey;not null" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`
}

type RolePermission struct {
	RoleID       string `gorm:"type:uuid;primaryKey;not null" json:"role_id"`
	PermissionID string `gorm:"type:uuid;primaryKey;not null" json:"permission_id"`
}
