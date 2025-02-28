package postgresql

import (
	"gorm.io/gorm"
)

func (p *Postgresql) SaveAllFields(model interface{}) (*gorm.DB, error) {
	result := p.Db.Save(model)
	if result.Error != nil {
		return result, result.Error
	}
	return result, nil
}

func (p *Postgresql) UpdateFields(model interface{}, updates interface{}, id string) (*gorm.DB, error) {
	result := p.Db.Model(model).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result, result.Error
	}
	return result, nil
}

func (p *Postgresql) SaveAllModelsFields(models []interface{}) (*gorm.DB, error) {
	// Use a transaction to ensure atomicity of updates
	tx := p.Db.Begin()
	if tx.Error != nil {
		return tx, tx.Error
	}

	// Loop through each model and update it
	for _, model := range models {
		result := tx.Save(model)
		if result.Error != nil {
			// If any update fails, rollback the transaction and return the error
			tx.Rollback()
			return result, result.Error
		}
	}

	// Commit the transaction if all updates are successful
	if err := tx.Commit().Error; err != nil {
		return tx, err
	}

	return tx, nil
}
