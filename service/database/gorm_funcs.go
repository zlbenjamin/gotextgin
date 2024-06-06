package database

// Get record by a primary key
func GetRecordByPk[T, PK any](st T, id PK) (int64, error) {
	db := GetDB()
	db2 := db.First(st, id)
	return db2.RowsAffected, db2.Error
}

// Delete one record by a primary key
func DeleteOneRecordByPk[T, PK any](st T, id PK) (int64, error) {
	db := GetDB()
	db2 := db.Delete(st, id)
	return db2.RowsAffected, db2.Error
}
