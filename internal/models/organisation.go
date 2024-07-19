package models

type Organisation struct {
	Orgid       string `json:"org_id" gorm:"primaryKey;unique;not null"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Users       []User `gorm:"many2many:user_organisations;foreignKey:Orgid;joinForeignKey:org_id;References:Userid;joinReferences:user_id" json:"users"`
}
