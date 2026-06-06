package models

import "gorm.io/gorm"

func MigrateTables(db *gorm.DB) {

	db.AutoMigrate(&Usuario{})
	db.AutoMigrate(&Evento{})

}
