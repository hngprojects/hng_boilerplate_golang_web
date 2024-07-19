package models

type Organisation struct {
	Orgid          string `gorm:"type:uuid;primaryKey;unique;not null" json:"org_id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Users       []User `gorm:"many2many:user_organisations;foreignKey:Orgid;joinForeignKey:org_id;References:Userid;joinReferences:user_id" json:"users"`
}
