package postgresql

func (p *Postgresql) DeleteRecordFromDb(record interface{}) error {
	tx := p.Db.Delete(record)
	return tx.Error
}

func (p *Postgresql) HardDeleteRecordFromDb(record interface{}) error {
	tx := p.Db.Unscoped().Delete(record)
	return tx.Error
}
