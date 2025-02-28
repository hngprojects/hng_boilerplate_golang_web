package templates

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func CreateTemplate(db *storage.Database, templateReq models.TemplateRequest) (models.EmailTemplate, error) {

	template := models.EmailTemplate{
		ID:   utility.GenerateUUID(),
		Name: templateReq.Name,
		Body: templateReq.Body,
	}

	err := template.Create(db.Postgresql)
	if err != nil {
		return template, err
	}
	return template, nil
}

func GetTemplates(db *gorm.DB) ([]models.EmailTemplate, error) {
	pdb := inst.InitDB(db)
	var template models.EmailTemplate

	templates, err := template.GetAllTemplates(pdb)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func GetTemplate(db *gorm.DB, id string) (models.EmailTemplate, error) {
	var template models.EmailTemplate
	pdb := inst.InitDB(db)

	temp, err := template.GetTemplateByID(pdb, id)
	if err != nil {
		return models.EmailTemplate{}, err
	}
	return temp, nil
}

func DeleteTemplate(db *gorm.DB, id string) error {
	var template models.EmailTemplate

	pdb := inst.InitDB(db)
	template.ID = id

	err := template.DeleteTemplate(pdb, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTemplate(db *gorm.DB, id string, templateReq models.EmailTemplate) (models.EmailTemplate, error) {
	var template models.EmailTemplate
	pdb := inst.InitDB(db)

	template.ID = id
	template.Name = templateReq.Name
	template.Body = templateReq.Body

	temp, err := template.UpdateTemplate(pdb, id)
	if err != nil {
		return temp, err
	}
	return temp, nil
}
