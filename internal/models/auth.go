package models

import (
	"errors"
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type PasswordReset struct {
	ID        string         `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	Email     string         `gorm:"index"`
	Token     string         `gorm:"uniqueIndex"`
	ExpiresAt time.Time      `gorm:"column:expires_at" json:"expires_at"`
	CreatedAt time.Time      `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type MagicLink struct {
	ID        string         `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	Email     string         `gorm:"index"`
	Token     string         `gorm:"uniqueIndex"`
	ExpiresAt time.Time      `gorm:"column:expires_at" json:"expires_at"`
	CreatedAt time.Time      `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type MagicLinkRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyMagicLinkRequest struct {
	Token string `json:"token" validate:"required"`
}

type ChangePasswordRequestModel struct {
	OldPassword string `json:"old_password" validate:""`
	NewPassword string `json:"new_password" validate:"required,min=7"`
}

type ForgotPasswordRequestModel struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequestModel struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=7"`
}

func (p *PasswordReset) CreatePasswordReset(db *gorm.DB) error {

	err := postgresql.CreateOneRecord(db, &p)

	if err != nil {
		return err
	}

	return nil
}

func (pr *PasswordReset) GetPasswordResetByToken(db *gorm.DB, token string) (PasswordReset, error) {
	var reset PasswordReset
	if err := db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&reset).Error; err != nil {
		return reset, err
	}
	return reset, nil
}

func (pr *PasswordReset) GetPasswordResetByEmail(db *gorm.DB, email string) (*PasswordReset, error) {
	var reset PasswordReset
	if err := db.Where("email = ?", email).First(&reset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reset, nil
}

func (pr *PasswordReset) DeletePasswordReset(db *gorm.DB) error {

	err := postgresql.DeleteRecordFromDb(db, pr)

	if err != nil {
		return err
	}

	return nil
}

func (m *MagicLink) CreateMagicLink(db *gorm.DB) error {

	err := postgresql.CreateOneRecord(db, &m)

	if err != nil {
		return err
	}

	return nil
}

func (m *MagicLink) GetMagicLinkByToken(db *gorm.DB, token string) (MagicLink, error) {
	var magic MagicLink
	if err := db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&magic).Error; err != nil {
		return magic, err
	}
	return magic, nil
}

func (m *MagicLink) GetMagicLinkByEmail(db *gorm.DB, email string) (*MagicLink, error) {
	var magic MagicLink
	if err := db.Where("email = ?", email).First(&magic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &magic, nil
}

func (m *MagicLink) DeleteMagicLink(db *gorm.DB) error {

	err := postgresql.DeleteRecordFromDb(db, m)

	if err != nil {
		return err
	}

	return nil
}
