package migrate

import (
	"app/database"
	model "app/modules/authen/model"
)

func MigrateAuthen() bool {
	db := database.DB

	db.AutoMigrate(&model.User{})

	return true
}
