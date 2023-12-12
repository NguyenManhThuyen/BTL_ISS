package model

import (
	"time"
	//"app/model"
	modelStudent "app/modules/student/model"
	modelAdvisor "app/modules/advisor/model"
	modelCouncil "app/modules/council/model"
	modelFacultyOffice "app/modules/facultyOffice/model"
	modelHeadOfSubject "app/modules/headOfSubject/model"
	"github.com/go-playground/validator/v10"
)

var StudentRole, AdvisorRole, HeadOfSubjectRole, FacultyOfficeRole, CouncilRole = 1, 2, 3, 4, 5

type SignUpInput struct {
	Id uint `json:"id"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,min=8"`
	Image           string `json:"image"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	FullName        string `json:"fullName"`
	Address         string `json:"address"`
	Gender          string `json:"gender"`
	PhoneNumber     string `json:"phoneNumber"`
	Birthday        string `json:"birthday"`
	Role            int    `json:"role"`
	Code 			string `json:"code"`
}

type SignInInput struct {
	Email string `json:"email" validate:"required"`
	Password string `json:"password"  validate:"required"`
	Role     int    `json:"role"`
}

type UserResponse struct {
	ID      uint    `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      int       `json:"role,omitempty"`
	Status bool `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FilterUserRecordStudent(user *modelStudent.Student) UserResponse {
	return UserResponse{
		ID:      user.ID,
		Email:     user.Email,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func FilterUserRecordAdvisor(user *modelAdvisor.Advisor) UserResponse {
	return UserResponse{
		ID:      user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func FilterUserRecordCouncil(user *modelCouncil.Council) UserResponse {
	return UserResponse{
		ID:      user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func FilterUserRecordFacultyOffice(user *modelFacultyOffice.FacultyOffice) UserResponse {
	return UserResponse{
		ID:      user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func FilterUserRecordHeadOfSubject(user *modelHeadOfSubject.HeadOfSubject) UserResponse {
	return UserResponse{
		ID:      user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
