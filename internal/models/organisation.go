package models

import "gorm.io/gorm"

type Organisation struct {
	OrgID       string `json:"org_id" gorm:"primaryKey;unique;not null"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Users       []User `gorm:"many2many:user_organisations;foreignKey:OrgID;joinForeignKey:org_id;References:UserID;joinReferences:user_id" json:"users"`
}

func AddUserToOrganisation(db *gorm.DB, orgID, userID string) error {
	// Add user to organisation
	err := db.Exec("INSERT INTO user_organisations (org_id, user_id) VALUES (?, ?)", orgID, userID).Error
	if err != nil {
		return err
	}
	return nil
}
