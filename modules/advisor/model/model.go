package model

import (
	"app/model"
	// "errors"
	// "time"
	// "gorm.io/gorm"
)

type Advisor struct {
	model.Info `gorm:"embedded;-:migration"`
	Code       string `json:"code" gorm:"column:CODE;size:10;not null"`
	ThesisID uint   `json:"thesisID" gorm:"column:THESIS_ID;index"`
}

type CreateAdvisor struct {
	ID          uint `json:"id"`
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

type UpdateAdvisor struct {
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

func (Advisor) TableName() string {
	return "tbl_advisor"
}

// func (s *Advisor) BeforeUpdate(tx *gorm.DB) (err error) {

// 	// Kiểm tra phiên bản cập nhật để cập nhật log (Log Version)
// 	if tx.Statement.Changed("LogVersion") {
// 		return errors.New("LOG_VERSION_CHANGED")
// 	}

// 	// Tăng log version lên 1 và cập nhật thời gian cập nhật
// 	s.LogVersion += 1
// 	tx.Statement.SetColumn("LogVersion", s.LogVersion)
// 	tx.Statement.SetColumn("UpdatedAt", time.Now())

// 	return nil
// }

// func (advisor *Advisor) BeforeCreate(tx *gorm.DB) (err error) {
// 	// Check if the username already exists
// 	var count int64
// 	tx.Model(&Advisor{}).Where("username = ?", advisor.Username).Count(&count)
// 	if count > 0 {
// 		return errors.New("CONFLICT_USERNAME") // You can customize the error based on your needs
// 	}

// 	tx.Model(&Advisor{}).Where("email = ?", advisor.Email).Count(&count)
// 	if count > 0 {
// 		return errors.New("CONFLICT_EMAIL") // You can customize the error based on your needs
// 	}

// 	// Set the CreatedAt timestamp
// 	advisor.CreatedAt = time.Now()

// 	return nil
// }
