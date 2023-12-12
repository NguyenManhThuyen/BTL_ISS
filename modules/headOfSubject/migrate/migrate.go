package headOfSubjectMigrate

import (
	"app/database"
	model "app/modules/headOfSubject/model"
)

func MigrateTable() bool {
	db := database.DB

	db.AutoMigrate(&model.HeadOfSubject{})

	return true
}