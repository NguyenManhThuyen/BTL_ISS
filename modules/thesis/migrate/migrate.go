package thesisRoute

import (
	"app/database"
	model "app/modules/thesis/model"
)

func MigrateTable() bool {
	db := database.DB

	db.AutoMigrate(&model.Thesis{})
	db.AutoMigrate(&model.ThesisTask{})
	db.AutoMigrate(&model.Mission{})
	db.AutoMigrate(&model.Program{})
	return true
}
