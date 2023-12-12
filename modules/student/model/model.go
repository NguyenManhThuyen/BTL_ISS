package model

import (
	"app/model"

)

type Student struct {
	model.Info
	Code        string `json:"code" gorm:"column:CODE;size:10"`
	Status      bool   `json:"status" gorm:"column:STATUS;default:false"`
	ThesisID uint   `json:"thesisID" gorm:"column:THESIS_ID;index"`
}

type CreateStudent struct {
	Id 			uint `json:"id"  validate:"required"`
	Code        string `json:"code"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	FullName    string `json:"fullName" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Birthday    string `json:"birthday" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Image       string `json:"image"`
}

type UpdateStudent struct {
	UUID        string `json:"uuid"`
	Code        string `json:"code"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phoneNumber"`
	Birthday    string `json:"birthday"`
	Password    string `json:"password"`
	Image       string `json:"image"`
	IsDeleted   bool   `json:"isDeleted"`
	Status      bool   `json:"status"`
}

// TableName overrides the table name used by GORM to tbl_student
func (Student) TableName() string {
	return "tbl_student"
}
