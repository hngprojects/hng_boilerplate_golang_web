package postgresql

import (
	"fmt"
)

func (p *Postgresql) CreateOneRecord(model interface{}) error {
	result := p.Db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("record creation failed")
	}
	return nil
}

func (p *Postgresql) CreateMultipleRecords(model interface{}, length int) error {
	result := p.Db.Create(model)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != int64(length) {
		return fmt.Errorf("record creation failed")
	}
	return nil
}
