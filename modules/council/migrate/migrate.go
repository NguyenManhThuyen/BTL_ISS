package councilMigrate

import (
	"app/database"
	model "app/modules/council/model"
)

func MigrateTable() bool {
	db := database.DB

	db.AutoMigrate(&model.Council{})

	return true
}