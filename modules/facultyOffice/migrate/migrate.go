package facultyOfficeMigrate

import (
	"app/database"
	model "app/modules/facultyOffice/model"
)

func MigrateTable() bool {
	db := database.DB

	db.AutoMigrate(&model.FacultyOffice{})

	return true
}
