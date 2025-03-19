package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseManager struct {
}

func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{}
}

func (d *DatabaseManager) GetConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db, nil
}

func (d *DatabaseManager) CloseConnection(conn *gorm.DB) {
	db, _ := conn.DB()
	db.Close()
}
