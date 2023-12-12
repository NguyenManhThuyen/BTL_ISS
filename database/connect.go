package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/cengsin/oracle"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func ConnectToOracle() bool {
	dsn := "system/manhthuyen@localhost:1521/ORCL"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,         // Disable color
		},
	)
	DB, _ = gorm.Open(oracle.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	return true
}
