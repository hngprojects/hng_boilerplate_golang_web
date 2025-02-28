package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"gorm.io/gorm"
)

type EmailTemplate struct {
	ID        string         `gorm:"type:uuid;primaryKey;unique; not null" json:"id"`
	Name      string         `gorm:"column:name; type:varchar(255); not null" json:"name"`
	Body      string         `gorm:"column:body; type:text; not null" json:"body"`
	CreatedAt time.Time      `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type TemplateRequest struct {
	Name string `json:"name" validate:"required"`
	Body string `json:"body" validate:"required"`
}

func (t *EmailTemplate) Create(db database.DatabaseManager) error {

	err := db.CreateOneRecord(&t)
	if err != nil {
		return err
	}
	return nil
}

func (t *EmailTemplate) GetAllTemplates(db database.DatabaseManager) ([]EmailTemplate, error) {
	var templates []EmailTemplate

	err := db.SelectAllFromDb("", "", &templates, "")
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (t *EmailTemplate) GetTemplateByID(db database.DatabaseManager, ID string) (EmailTemplate, error) {
	var template EmailTemplate

	exists := db.CheckExists(&template, "id = ?", ID)
	if !exists {
		return EmailTemplate{}, gorm.ErrRecordNotFound
	}

	nerr, err := db.SelectOneFromDb(&template, "id = ?", t.ID)
	if err != nil {
		return EmailTemplate{}, nerr
	}
	return template, nil
}

func (t *EmailTemplate) DeleteTemplate(db database.DatabaseManager, ID string) error {

	exists := db.CheckExists(t, "id = ?", ID)
	if !exists {
		return gorm.ErrRecordNotFound
	}

	err := db.DeleteRecordFromDb(&t)
	if err != nil {
		return err
	}
	return nil
}

func (t *EmailTemplate) UpdateTemplate(db database.DatabaseManager, ID string) (EmailTemplate, error) {
	t.ID = ID

	exists := db.CheckExists(&EmailTemplate{}, "id = ?", ID)
	if !exists {
		return EmailTemplate{}, gorm.ErrRecordNotFound
	}

	_, err := db.SaveAllFields(t)
	if err != nil {
		return EmailTemplate{}, err
	}

	updatedTemplate := EmailTemplate{}
	err = db.DB().First(&updatedTemplate, "id = ?", ID).Error
	if err != nil {
		return EmailTemplate{}, err
	}

	return updatedTemplate, nil
}
