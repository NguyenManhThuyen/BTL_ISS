package controller

import (
	"app/config"
	"app/database"
	"encoding/json"

	modelll "app/modules/advisor/model"
	modell "app/modules/student/model"
	"app/modules/thesis/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ThesisStatusResponse struct {
	Status bool `json:"status"`
}

type Input struct {
	ThesisUUID  string `json:"thesisUUID"`
	StudentID string `json:"StudentID"`
}

func generateRandomUUID() string {
	return uuid.NewString()
}

// @title Student API
// @version 1.0
// @description API for managing student data
// @termsOfService http://swagger.io/terms/
// @BasePath /student
// @schemes http
// @produce json
// @consumes json

// UpdateThesisApprovalStatus updates the ApprovalStatus of a thesis
// @Summary Update the ApprovalStatus of a thesis
// @Description Update the ApprovalStatus of a thesis
// @Tags Thesis
// @Accept json
// @Produce json
// @Param body body model.ApprovalStatusForThesis true "ApprovalStatus information"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis/approval [put]
func UpdateThesisApprovalStatus(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var payload model.ApprovalStatusForThesis
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = "Failed to parse request body"
		return c.JSON(response)
	}

	var thesis model.Thesis
	if err := db.First(&thesis, payload.ThesisID).Error; err != nil {
		response.Status = false
		response.Message = "Thesis not found"
		return c.JSON(response)
	}

	// Cập nhật giá trị ApprovalStatus của thesis
	if err := db.Model(&thesis).Update("APPROVAL_STATUS", payload.ApprovalStatus).Error; err != nil {
		response.Status = false
		response.Message = "Update error"
		return c.JSON(response)
	}
	response.Data= thesis
	response.Status = true
	response.Message = "Thesis ApprovalStatus updated successfully"
	return c.JSON(response)
}


// GetThesesByCreateBy gets the theses created by a user based on UUID
// @Summary Get theses created by a user
// @Description Get theses created by a user based on UUID
// @Tags Thesis
// @Produce json
// @Param createBy path string true "UUID of the user who created the theses"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis/get-by-createby/{createBy} [get]
func GetThesesByCreateBy(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	// Get the CreateBy UUID from the request parameters
	createByUUID := c.Params("createBy")

	var theses []model.Thesis
	if err := db.Preload("Missions").Preload("Programs").Preload("ThesisTask").Preload("Students").Preload("Advisors").
		Where("created_by = ?", createByUUID).Find(&theses).Error; err != nil {
		response.Status = false
		response.Message = "Failed to fetch theses"
		return c.JSON(response)
	}

	response.Status = true
	response.Data = theses
	return c.JSON(response)
}

// ListTheses lấy danh sách tất cả luận văn
// @Summary Get a list of theses
// @Description Trả về danh sách tất cả các luận văn với thông tin chi tiết.
// @Tags Thesis
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis [get]
func ListTheses(c *fiber.Ctx) error {
    response := new(config.DataResponse)
    db := database.DB

    var theses []model.Thesis
    if err := db.Preload("Missions").Preload("Programs").Preload("ThesisTask").Preload("Students").Preload("Advisors").Find(&theses).Error; err != nil {
        response.Status = false
        response.Message = "Failed to fetch theses"
        return c.JSON(response)
    }

    response.Status = true
    response.Data = theses
    return c.JSON(response)
}


// GetThesis trả về danh sách tất cả Thesis
// @Summary Get a list of theses
// @Description Trả về danh sách luận văn.
// @Tags Thesis
// @Produce json
// @Param uuid path string true "Thesis UUID"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis/{uuid} [get]
func GetThesis(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var thesis model.Thesis
	if err := db.Preload("Missions").Preload("Programs").Preload("ThesisTask").Preload("Students").Preload("Advisors").First(&thesis, "uuid = ?", c.Params("uuid")).Error; err != nil {
		response.Status = false
		response.Message = "Thesis not found"
		return c.JSON(response)
	}

	response.Status = true
	response.Data = thesis
	return c.JSON(response)
}

// PostStatusThesis trả về trạng thái của một Thesis cho một Student cụ thể.
// @Summary Get the status of a Thesis for a specific Student
// @Description Trả về trạng thái của một Thesis cho một Student cụ thể.
// @Tags Thesis
// @Produce json
// @Param body body Input true "Thesis information"
// @Success 200 {object} config.DataResponse
// @Failure 400 {object} config.DataResponse
// @Failure 404 {object} config.DataResponse
// @Router /thesis/status-thesis [post]
func PostStatusThesis(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var payload = Input{}
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = "Failed to parse request body"
		return c.JSON(response)
	}

	var thesis model.Thesis
	if err := db.Preload("Missions").Preload("Programs").Preload("ThesisTask").Preload("Students").Preload("Advisors").First(&thesis, "uuid = ?", payload.ThesisUUID).Error; err != nil {
		response.Status = false
		response.Message = "Thesis not found"
		return c.JSON(response)
	}

	var student modell.Student
	results := database.DB.First(&student,payload.StudentID)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}


	// Tạo một bản sao của thesis
	response.Status = true
	return c.JSON(response)
}

// CreateThesis creates a new Thesis
// @Summary Create a new thesis
// @Description Create a new thesis
// @Tags Thesis
// @Accept json
// @Produce json
// @Param body body []model.CreateThesis true "Thesis information"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis [post]
func CreateThesis(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var payload []model.CreateThesis
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = "Failed to parse request body"
		return c.JSON(response)
	}

	tx := db.Begin()
	defer tx.Commit()

	for _, thesisPayload := range payload {
		newThesis := model.Thesis{
			TitleVi:        thesisPayload.TitleVi,
			TitleEn:        thesisPayload.TitleEn,
			ApprovalStatus: thesisPayload.ApprovalStatus,
			ThesisType:     thesisPayload.ThesisType,
			Semester:       thesisPayload.Semester,
			UserRoleOwner:  thesisPayload.UserRoleOwner,
			ThesisInfo:     thesisPayload.ThesisInfo,
			StartTime:      thesisPayload.StartTime,
			EndTime:        thesisPayload.EndTime,
		}
		newThesis.ID = thesisPayload.ID
		// Handle Advisors
		// Create the thesis
		if err := db.Create(&newThesis).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = "Failed to create thesis"
			return c.JSON(response)
		}

		for _, advisorIDPayload := range thesisPayload.Advisors {
			var advisor modelll.Advisor
			if err := db.First(&advisor, advisorIDPayload.AdvisorID).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = "Advisor not found"
				return c.JSON(response)
			}
		
			// Cập nhật giá trị của cột THESIS_ID
			if err := db.Model(&advisor).UpdateColumn("THESIS_ID", newThesis.ID).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = "Update error"
				return c.JSON(response)
			}
		}
		
		/// ???
		for _, studentIDPayload := range thesisPayload.Students {

			var student modell.Student
			if err := db.First(&student, studentIDPayload.StudentID).Error; err != nil {
				response.Status = false
				response.Message = "Student not found"
				return c.JSON(response)
			}
			if err := db.Model(&student).UpdateColumn("THESIS_ID", newThesis.ID).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = "Update error"
				return c.JSON(response)
			}
		}
		// Handle Missions
		for _, missionPayload := range thesisPayload.Missions {
			newMission := model.Mission{
				Value:    missionPayload.Value,
				ThesisID: newThesis.ID,
			}
			newMission.ID = missionPayload.ID
			if err := db.Create(&newMission).Error; err != nil {
				response.Status = false
				response.Message = "Failed to create mission"
				return c.JSON(response)
			}
		}

		// Handle Programs
		for _, programPayload := range thesisPayload.Programs {
			newProgram := model.Program{
				Value:    programPayload.Value,
				ThesisID: newThesis.ID,
			}
			newProgram.ID = programPayload.ID
			if err := db.Create(&newProgram).Error; err != nil {
				response.Status = false
				response.Message = "Failed to create program"
				return c.JSON(response)
			}
		}

		// Create thesis tasks for this thesis
		for _, taskPayload := range thesisPayload.ThesisTask {
			newTask := model.ThesisTask{
				Title:       taskPayload.Title,
				Deadline:    taskPayload.Deadline,
				Status:      taskPayload.Status,
				Priority:    taskPayload.Priority,
				Description: taskPayload.Description,
				Note:        taskPayload.Note,
				StartTime:   taskPayload.StartTime,
				EndTime:     taskPayload.EndTime,
				ThesisID:    newThesis.ID,
			}
			newTask.ID = taskPayload.ID
			if err := db.Create(&newTask).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = "Failed to create thesis task"
				return c.JSON(response)
			}
		}
	}

	response.Status = true
	response.Message = "Thesis(s) created successfully"
	return c.JSON(response)
}

// UpdateThesis updates the information of a Thesis based on UUID
// @Summary Update thesis details by UUID
// @Description Update thesis details by UUID
// @Tags Thesis
// @Accept json
// @Produce json
// @Param body body []model.UpdateThesis true "Thesis information to update"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis [put]
func UpdateThesis(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var payload []model.UpdateThesis
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = "Failed to parse request body"
		return c.JSON(response)
	}

	tx := db.Begin()
	defer tx.Commit()

	for _, thesisPayload := range payload {
		// Find the thesis by ID
		var thesis model.Thesis
		if err := db.Preload("Missions").Preload("Programs").Preload("ThesisTask").Preload("Students").Preload("Advisors").First(&thesis, "uuid = ?", thesisPayload.UUID).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = "Thesis not found"
			return c.JSON(response)
		}

		// Clear the old lists of students, advisors, missions, and programs
		// Xóa dữ liệu cũ
		db.Model(&thesis).Association("Students").Clear()
		db.Model(&thesis).Association("Advisors").Clear()
		db.Delete(&thesis.Missions)
		db.Delete(&thesis.Programs)
		db.Delete(&thesis.ThesisTask)

		// Handle Students
		for _, studentIDPayload := range thesisPayload.Students {
			var student modell.Student
			if err := db.First(&student, "uuid = ?", studentIDPayload.StudentID).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = "Student not found"
				return c.JSON(response)
			}
			thesis.Students = append(thesis.Students, student)
		}

		// Handle Advisors
		for _, advisorIDPayload := range thesisPayload.Advisors {
			var advisor modelll.Advisor
			if err := db.First(&advisor, "uuid = ?", advisorIDPayload.AdvisorID).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = "Advisor not found"
				return c.JSON(response)
			}
			thesis.Advisors = append(thesis.Advisors, advisor)
		}

		// Handle Missions
		for _, missionPayload := range thesisPayload.Missions {
			newMission := model.Mission{
				Value:    missionPayload.Value,
				ThesisID: thesis.ID,
			}

			if err := db.Create(&newMission).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = "Failed to create mission"
				return c.JSON(response)
			}
		}

		// Handle Programs
		for _, programPayload := range thesisPayload.Programs {
			newProgram := model.Program{
				Value:    programPayload.Value,
				ThesisID: thesis.ID,
			}

			if err := db.Create(&newProgram).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = "Failed to create program"
				return c.JSON(response)
			}
		}

		// Update thesis information if values are present in thesisPayload
		if thesis.TitleVi != thesisPayload.TitleVi {
			thesis.TitleVi = thesisPayload.TitleVi
		}
		if thesis.TitleEn != thesisPayload.TitleEn {
			thesis.TitleEn = thesisPayload.TitleEn
		}
		if thesis.ApprovalStatus != thesisPayload.ApprovalStatus {
			thesis.ApprovalStatus = thesisPayload.ApprovalStatus
		}
		if thesis.ThesisType != thesisPayload.ThesisType {
			thesis.ThesisType = thesisPayload.ThesisType
		}
		if thesis.Semester != thesisPayload.Semester {
			thesis.Semester = thesisPayload.Semester
		}
		if thesis.UserRoleOwner != thesisPayload.UserRoleOwner {
			thesis.UserRoleOwner = thesisPayload.UserRoleOwner
		}
		if thesis.ThesisInfo != thesisPayload.ThesisInfo {
			thesis.ThesisInfo = thesisPayload.ThesisInfo
		}
		if !thesis.StartTime.Equal(thesisPayload.StartTime) {
			thesis.StartTime = thesisPayload.StartTime
		}
		if !thesis.EndTime.Equal(thesisPayload.EndTime) {
			thesis.EndTime = thesisPayload.EndTime
		}

		// Update thesis tasks
		for _, taskPayload := range thesisPayload.ThesisTask {
			var task model.ThesisTask
			if taskPayload.UUID != "" {
				if err := db.First(&task, "uuid = ?", taskPayload.UUID).Error; err != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("NOT_ID_EXISTS")
					return c.JSON(response)
				}
			}
			task.Deadline = taskPayload.Deadline
			task.Description = taskPayload.Description
			task.EndTime = taskPayload.EndTime
			task.Note = taskPayload.Note
			task.Priority = taskPayload.Priority
			task.StartTime = taskPayload.StartTime
			task.Status = taskPayload.Status
			task.Title = taskPayload.Title
			task.ThesisID = thesis.ID

			if err := db.Save(&task).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = config.GetMessageCode("SYSTEM_ERROR")
				return c.JSON(response)
			}
		}

		// Save the updated thesis
		if err := db.Save(&thesis).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = "Failed to update thesis"
			return c.JSON(response)
		}

		for _, studentIDPayload := range thesisPayload.Students {
			var student modell.Student
			if err := tx.First(&student, "uuid = ?", studentIDPayload.StudentID).Error; err != nil {
				response.Status = false
				response.Message = "Student not found"
				return c.JSON(response)
			}
		}
	}

	response.Status = true
	response.Message = "Thesis(s) updated successfully"
	return c.JSON(response)
}

// DeleteThesis xóa một Thesis dựa trên ID
// @Summary Xóa Thesis
// @Description Xóa một Thesis dựa trên UUID
// @Tags Thesis
// @Produce json
// @Param id path string true "ID của Thesis"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis/{uuid} [delete]
func DeleteThesis(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the thesis ID from the request parameters
	thesisUUID := c.Params("uuid")

	// Start a database transaction
	tx := database.DB.Begin()
	defer tx.Commit()

	// Find the thesis by UUID
	var thesis model.Thesis
	if err := tx.First(&thesis, "uuid = ?", thesisUUID).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Thesis not found"
		return c.JSON(response)
	}

	// Delete the thesis and its associated tasks
	if err := tx.Delete(&thesis).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Failed to delete thesis"
		return c.JSON(response)
	}

	response.Status = true
	response.Message = "Thesis deleted successfully"
	return c.JSON(response)
}

// AddStudentToThesis thêm một Student vào một Thesis dựa trên UUID
// @Summary Add a student to a thesis
// @Description Add a student to a thesis by UUID
// @Tags Thesis
// @Produce json
// @Param thesisUUID path string true "UUID của Thesis"
// @Param StudentID path string true "UUID của Student"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis/addstudent/{thesisUUID}/{StudentID} [post]
func AddStudentToThesis(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	thesisUUID := c.Params("thesisUUID")
	StudentID := c.Params("StudentID")

	tx := database.DB.Begin()
	defer tx.Commit()

	var thesis model.Thesis

	if err := tx.First(&thesis, "uuid = ? ", thesisUUID).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Thesis not found"
		return c.JSON(response)
	}

	var student modell.Student
	if err := tx.First(&student, "uuid = ? ", StudentID).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Student not found"
		return c.JSON(response)
	}

	if err := tx.Model(&thesis).Association("Students").Append(&student); err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Failed to add student to thesis"
		return c.JSON(response)
	}

	response.Status = true
	response.Message = "Student added to thesis successfully"
	return c.JSON(response)
}

// AddAdvisorToThesis thêm một Advisor vào một Thesis dựa trên UUID
// @Summary Add an advisor to a thesis
// @Description Add an advisor to a thesis by UUID
// @Tags Thesis
// @Produce json
// @Param thesisUUID path string true "UUID của Thesis"
// @Param advisorUUID path string true "UUID của Advisor"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis/addadvisor/{thesisUUID}/{advisorUUID} [post]
func AddAdvisorToThesis(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	//db := database.DB.Begin()

	// Get the thesis ID and advisor ID from the request parameters
	thesisUUID := c.Params("thesisUUID")
	advisorUUID := c.Params("advisorUUID")

	tx := database.DB.Begin()
	defer tx.Commit()

	// Find the thesis by its UUID
	var thesis model.Thesis
	if err := tx.First(&thesis, "uuid = ? ", thesisUUID).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Thesis not found"
		return c.JSON(response)
	}

	// Find the advisor by its UUID
	var advisor modelll.Advisor
	if err := tx.First(&advisor, "uuid = ? ", advisorUUID).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Advisor not found"
		return c.JSON(response)
	}

	// Add the advisor to the thesis
	if err := tx.Model(&thesis).Association("Advisors").Append(&advisor).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Failed to add advisor to thesis"
		return c.JSON(response)
	}

	response.Status = true
	response.Message = "Advisor added to thesis successfully"
	return c.JSON(response)
}

// RemoveStudentFromThesis xóa một Student khỏi một Thesis dựa trên UUID
// @Summary Remove a student from a thesis
// @Description Remove a student from a thesis by UUID
// @Tags Thesis
// @Produce json
// @Param thesisUUID path string true "UUID của Thesis"
// @Param StudentID path string true "UUID của Student"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis/removestudent/{thesisUUID}/{StudentID} [delete]
func RemoveStudentFromThesis(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	thesisUUID := c.Params("thesisUUID")
	StudentID := c.Params("StudentID")

	tx := database.DB.Begin()
	defer tx.Commit()

	var thesis model.Thesis

	if err := tx.First(&thesis, "uuid = ? ", thesisUUID).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Thesis not found"
		return c.JSON(response)
	}

	var student modell.Student
	if err := tx.First(&student, "uuid = ? ", StudentID).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Student not found"
		return c.JSON(response)
	}

	if err := tx.Model(&thesis).Association("Students").Delete(&student); err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Failed to remove student from thesis"
		return c.JSON(response)
	}

	response.Status = true
	response.Message = "Student removed from thesis successfully"
	return c.JSON(response)
}

// RemoveAdvisorFromThesis xóa một Advisor khỏi một Thesis dựa trên UUID
// @Summary Remove an advisor from a thesis
// @Description Remove an advisor from a thesis by UUID
// @Tags Thesis
// @Produce json
// @Param thesisUUID path string true "UUID của Thesis"
// @Param advisorUUID path string true "UUID của Advisor"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis/removeadvisor/{thesisUUID}/{advisorUUID} [delete]
func RemoveAdvisorFromThesis(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the thesis UUID and advisor UUID from the request parameters
	thesisUUID := c.Params("thesisUUID")
	advisorUUID := c.Params("advisorUUID")

	tx := database.DB.Begin()
	defer tx.Commit()

	// Find the thesis and advisor by their UUIDs
	var thesis model.Thesis
	var advisor modelll.Advisor

	if err := tx.First(&thesis, "uuid = ? ", thesisUUID).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Thesis not found"
		return c.JSON(response)
	}

	if err := tx.First(&advisor, "uuid = ? ", advisorUUID).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Advisor not found"
		return c.JSON(response)
	}

	// Remove the advisor from the thesis
	if err := tx.Model(&thesis).Association("Advisors").Delete(&advisor).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Failed to remove advisor from thesis"
		return c.JSON(response)
	}

	response.Status = true
	response.Message = "Advisor removed from thesis successfully"
	return c.JSON(response)
}

// CreateTestTheses tạo dữ liệu luận văn thử nghiệm
// @Summary Create test thesis data
// @Description Create test thesis data with real data for development or testing purposes
// @Tags Thesis
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /thesis/create-test [post]
func CreateTestTheses(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	tx := database.DB.Begin()

	// Mảng JSON chứa dữ liệu luận văn thử nghiệm
	testThesesJSON := []byte(`[
    {
    "title_vi": "Đề tài Tiếng Việt 3",
    "title_en": "Thesis English 3",
    "approval_status": 1,
    "thesis_type": 1,
    "semester": "2024A",
    "program": [
      {"value": 1}
    ],
    "user_role_owner": 3,
    "thesis_info": "Thông tin đề tài 3",
    "thesis_task": [
      {
        "title": "Task 3",
        "deadline": "2024-01-10",
        "status": "Completed",
        "priority": 3,
        "description": "Task description",
        "note": "Task note",
        "start_time": "2023-12-20T09:00:00Z",
        "end_time": "2023-12-30T16:00:00Z"
      }
    ],
    "students": [
    ],
    "advisors": [
    ],
    "missions": [
    ]
  }
  
  ]`)

	var payload []model.CreateThesis

	if err := json.Unmarshal(testThesesJSON, &payload); err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = "Failed to unmarshal test thesis data"
		return c.JSON(response)
	}

	for _, thesisPayload := range payload {
		// Create a new thesis for each item in the payload
		newThesis := model.Thesis{
			TitleVi:        thesisPayload.TitleVi,
			TitleEn:        thesisPayload.TitleEn,
			ApprovalStatus: thesisPayload.ApprovalStatus,
			ThesisType:     thesisPayload.ThesisType,
			Semester:       thesisPayload.Semester,
			// Program:        thesisPayload.Program,
			UserRoleOwner: thesisPayload.UserRoleOwner,
			ThesisInfo:    thesisPayload.ThesisInfo,
		}

		if err := tx.Create(&newThesis).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = "Failed to create thesis"
			return c.JSON(response)
		}

		// Create thesis tasks for this thesis
		for _, taskPayload := range thesisPayload.ThesisTask {
			newTask := model.ThesisTask{
				Title:       taskPayload.Title,
				Deadline:    taskPayload.Deadline,
				Status:      taskPayload.Status,
				Priority:    taskPayload.Priority,
				Description: taskPayload.Description,
				Note:        taskPayload.Note,
				StartTime:   taskPayload.StartTime,
				EndTime:     taskPayload.EndTime,
				ThesisID:    newThesis.ID,
			}

			if err := tx.Create(&newTask).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = "Failed to create thesis task"
				return c.JSON(response)
			}
		}
	}

	tx.Commit()
	response.Status = true
	response.Message = "Thesis(s) created successfully"
	return c.JSON(response)
}
