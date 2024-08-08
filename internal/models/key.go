package models

type Key struct {
	ID     string `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	UserID string `gorm:"type:uuid;not null" json:"user_id"`
	Key    string `gorm:"type:varchar(255);not null" json:"key"`
	User   *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

type VerifyKeyRequestModel struct {
	Key string `json:"key" validate:"required"`
}
