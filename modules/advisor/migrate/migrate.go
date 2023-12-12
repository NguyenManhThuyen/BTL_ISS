package advisorMigrate

import (
	"app/database"
	model "app/modules/advisor/model"
)

func MigrateTable() bool {
	db := database.DB

	db.AutoMigrate(&model.Advisor{})

	return true
}