package model

import (
	"app/model"
)

type Council struct {
	model.Info `gorm:"embedded;-:migration"`
	Code       string `json:"code" gorm:"column:CODE;size:10;not null"`
}

type CreateCouncil struct {
	Id          uint   `json:"id"`
	Code        string `json:"code" validate:"required"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	FullName    string `json:"fullName" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Birthday    string `json:"birthday" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Image       string `json:"image" validate:"required"`
}

type UpdateCouncil struct {
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
}

func (Council) TableName() string {
	return "tbl_council"
}
