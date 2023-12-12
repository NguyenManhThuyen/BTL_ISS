package model

import (
    "app/model"
    modelll "app/modules/advisor/model"
    modell "app/modules/student/model"
    "time"
)


type Program struct {
	model.Header
	Value    int  `json:"value" validate:"required" gorm:"column:VALUE"`
	ThesisID uint `json:"thesisID" gorm:"column:THESIS_ID;index"`
}

type Mission struct {
	model.Header
	Value    string `json:"value" validate:"required" gorm:"column:VALUE"`
	ThesisID uint   `json:"thesisID" gorm:"column:THESIS_ID;index"`
}

type CreateProgram struct {
	ID uint `json:"id"`
	Value int `json:"value" validate:"required"`
}

type UpdateProgram struct {
	ID    uint `json:"id"`
	Value int    `json:"value" validate:"required"`
}

type UpdateMission struct {
	ID    uint   `json:"id"`
	Value string `json:"value" validate:"required"`
}

type CreateMission struct {
	ID uint `json:"id"`
	Value string `json:"value" validate:"required"`
}

type Thesis struct {
	model.Header
	TitleVi        string            `json:"titleVi" validate:"required" gorm:"column:TITLE_VI"`
	TitleEn        string            `json:"titleEn" validate:"required" gorm:"column:TITLE_EN"`
	ApprovalStatus int               `json:"approvalStatus" validate:"required" gorm:"column:APPROVAL_STATUS"`
	ThesisType     int               `json:"thesisType" validate:"required" gorm:"column:THESIS_TYPE"`
	Semester       string            `json:"semester" validate:"required" gorm:"column:SEMESTER"`
	UserRoleOwner int               `json:"userRoleOwner" gorm:"column:USER_ROLE_OWNER"`
	ThesisInfo    string            `json:"thesisInfo" gorm:"column:THESIS_INFO"`
	ThesisTask    []ThesisTask      `json:"thesisTask" gorm:"foreignKey:THESIS_ID"`
	Students      []modell.Student  `json:"students"  gorm:"foreignKey:THESIS_ID"`
	Advisors      []modelll.Advisor `json:"advisors"  gorm:"foreignKey:THESIS_ID"`
	Missions      []Mission         `json:"missions" gorm:"foreignKey:THESIS_ID"`
	Programs      []Program         `json:"programs" gorm:"foreignKey:THESIS_ID"`
	StartTime time.Time `json:"startTime" gorm:"column:START_TIME"`
	EndTime   time.Time `json:"endTime" gorm:"column:END_TIME"`
}

type ThesisTask struct {
	model.Header
	Title       string    `json:"title" validate:"required" gorm:"column:TITLE"`
	Deadline    string    `json:"deadline" validate:"required" gorm:"column:DEADLINE"`
	Status      string    `json:"status" validate:"required" gorm:"column:STATUS"`
	Priority    int       `json:"priority" validate:"required" gorm:"column:PRIORITY"`
	Description string    `json:"description" gorm:"column:DESCRIPTION"`
	Note        string    `json:"note" gorm:"column:NOTE"`
	ThesisID    uint      `json:"thesisID" gorm:"column:THESIS_ID;index"`
	StartTime   time.Time `json:"startTime" gorm:"column:START_TIME"`
	EndTime     time.Time `json:"endTime" gorm:"column:END_TIME"`
}

type CreateThesis struct {
	ID uint `json:"id"`
	TitleVi        string                    `json:"titleVi" validate:"required"`
	TitleEn        string                    `json:"titleEn" validate:"required"`
	ApprovalStatus int                       `json:"approvalStatus" validate:"required"`
	ThesisType     int                       `json:"thesisType" validate:"required"`
	Semester       string                    `json:"semester" validate:"required"`
	Programs       []CreateProgram           `json:"programs"`
	UserRoleOwner  int                       `json:"userRoleOwner"`
	ThesisInfo     string                    `json:"thesisInfo"`
	ThesisTask     []CreateThesisTask        `json:"thesisTask"`
	Students       []CreateStudentForThesis `json:"students"`
	Advisors       []CreateAdvisorForThesis `json:"advisors"`
	Missions       []CreateMission           `json:"missions"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type UpdateThesis struct {
	UUID           string                    `json:"uuid"`
	TitleVi        string                    `json:"titleVi"`
	TitleEn        string                    `json:"titleEn"`
	ApprovalStatus int                       `json:"approvalStatus"`
	ThesisType     int                       `json:"thesisType"`
	Semester       string                    `json:"semester"`
	Programs       []UpdateProgram           `json:"programs"`
	UserRoleOwner  int                       `json:"userRoleOwner"`
	ThesisInfo     string                    `json:"thesisInfo"`
	ThesisTask     []UpdateThesisTask        `json:"thesisTask"`
	Students       []*CreateStudentForThesis `json:"students"`
	Advisors       []*CreateAdvisorForThesis `json:"advisors"`
	Missions       []UpdateMission           `json:"missions"`
	StartTime      time.Time                 `json:"startTime"`
	EndTime        time.Time                 `json:"endTime"`
}

type CreateThesisTask struct {
	ID uint `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Deadline    string    `json:"deadline" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	Priority    int       `json:"priority" validate:"required"`
	Description string    `json:"description"`
	Note        string    `json:"note"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}

type UpdateThesisTask struct {
	UUID        string    `json:"uuid"`
	Title       string    `json:"title" validate:"required"`
	Deadline    string    `json:"deadline" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	Priority    int       `json:"priority" validate:"required"`
	Description string    `json:"description"`
	Note        string    `json:"note"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}

type CreateStudentForThesis struct {
	StudentID uint `json:"id"`
}

type CreateAdvisorForThesis struct {
	AdvisorID uint `json:"id"`
}

type ApprovalStatusForThesis struct {
	ApprovalStatus int
	ThesisID uint `json:"thesis_id"` 
}

func (Thesis) TableName() string {
	return "TBL_THESIS"
}

func (ThesisTask) TableName() string {
	return "TBL_THESIS_TASK"
}

func (Mission) TableName() string {
	return "TBL_THESIS_MISSIONS"
}

func (Program) TableName() string {
	return "TBL_THESIS_PROGRAMS"
}
