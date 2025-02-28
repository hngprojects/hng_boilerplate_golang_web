package postgresql

func (p *Postgresql) CountRecords(model interface{}) (int64, error) {
	var count int64
	result := p.Db.Model(model).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (p *Postgresql) CountSpecificRecords(model interface{}, query interface{}) (int64, error) {
	var count int64
	result := p.Db.Model(model).Where(query).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
