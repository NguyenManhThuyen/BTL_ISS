package model

import (
	"time"
	"gorm.io/gorm"
)

type Header struct {
	ID         uint      `gorm:"primaryKey;column:ID" `
	CreatedBy  string    `gorm:"column:CREATED_BY;size:50"`
	UpdatedBy  string    `gorm:"column:UPDATED_BY;size:50"`
	DeletedBy  string    `gorm:"column:DELETED_BY;size:50"`
	CreatedAt  time.Time `gorm:"column:CREATED_AT"`
	UpdatedAt  time.Time `gorm:"column:UPDATED_AT"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	LogVersion int64     `gorm:"column:LOG_VERSION;default:0" json:"-"`
	IsDeleted  bool      `gorm:"column:IS_DELETED;default:false"`
}

type Info struct {
	Header
	Password    string `json:"-"  validate:"required" gorm:"column:PASSWORD"`
	Image       string `json:"IMAGE"                        gorm:"column:IMAGE"`
	FirstName   string `json:"FIRST_NAME"                    gorm:"column:FIRST_NAME"`
	LastName    string `json:"LAST_NAME"                     gorm:"column:LAST_NAME"`
	FullName    string `json:"FULL_NAME"                     gorm:"column:FULL_NAME"`
	Email       string `json:"EMAIL" validate:"required"    gorm:"column:EMAIL"`
	Address     string `json:"ADDRESS"                      gorm:"column:ADDRESS"`
	Gender      string `json:"GENDER"                          gorm:"column:GENDER"`
	PhoneNumber string `json:"PHONE_NUMBER"                  gorm:"column:PHONE_NUMBER"`
	Birthday    string `json:"birthday" validate:"required" gorm:"column:BIRTHDAY"`
	Role        int    `json:"ROLE" validate:"required"      gorm:"column:ROLE"`
}
