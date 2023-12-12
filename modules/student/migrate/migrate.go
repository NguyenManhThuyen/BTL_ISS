package studentMigrate

import (
	"app/database"
	model "app/modules/student/model"
)

func MigrateTable() bool {
	db := database.DB

	db.AutoMigrate(&model.Student{})

	return true
}