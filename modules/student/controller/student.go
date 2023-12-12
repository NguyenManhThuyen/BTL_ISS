package controller

import (
	"app/config"
	"app/controller"
	"app/database"

	"app/modules/student/model"
	modelUsers "app/modules/users/model"
	"encoding/json"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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

// GetStudent trả về danh sách tất cả Student
// @Summary Get a list of students
// @Description Trả về danh sách sinh viên.
// @Tags Student
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /student [get]
func GetStudent(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var students []model.Student
	results := database.DB.Select("*").Order("id").Find(&students)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = students
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetStudentByUUID trả về thông tin của một Student dựa trên UUID
// @Summary Get student by UUID
// @Description Get student details by UUID
// @Tags Student
// @Produce json
// @Param uuid path string true "Student UUID"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /student/{uuid} [get]
func GetStudentByUUID(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var student model.Student
	results := database.DB.First(&student, "uuid = ?", c.Params("uuid"))
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = student
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetStudentByCODE trả về thông tin của một Student dựa trên CODE
// @Summary Get student by CODE
// @Description Get student details by CODE
// @Tags Student
// @Produce json
// @Param code path string true "Student CODE"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /student/code/{code} [get]
func GetStudentByCode(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var student model.Student
	results := database.DB.First(&student, "code = ?", c.Params("code"))
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = student
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// CreateStudent tạo một Student mới
// @Summary Create a new student
// @Description Create a new student
// @Tags Student
// @Accept json
// @Produce json
// @Param body body []model.CreateStudent true "Student information"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /student [post]
func CreateStudent(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	var payload []*model.CreateStudent
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := database.DB.Begin()
	defer tx.Commit()

	for _, item := range payload {
		listCheck := []string{"Code", "FirstName", "LastName", "FullName", "Email", "PhoneNumber", "Birthday", "Address", "Gender"}
		vItem := map[string]string{
			"Code":        item.Code,
			"FirstName":   item.FirstName,
			"LastName":    item.LastName,
			"FullName":    item.FullName,
			"Email":       item.Email,
			"PhoneNumber": item.PhoneNumber,
			"Birthday":    item.Birthday,
			"Address":     item.Address,
			"Gender":      item.Gender,
		}

		// Kiểm tra tất cả trường cùng một lúc
		missingFields := []string{}
		for _, field := range listCheck {
			if vItem[field] == "" {
				missingFields = append(missingFields, field)
			}
		}

		if len(missingFields) > 0 {
			tx.Rollback()
			response.Status = false
			response.Message = config.GetMessageCode("MISSING_FIELDS")
			response.ValidateError = missingFields
			return c.JSON(response)
		}

		newUser := new(model.Student)
		newUser.ID = item.Id
		newUser.Code = item.Code
		newUser.FirstName = item.FirstName
		newUser.LastName = item.LastName
		newUser.Image = item.Image
		newUser.PhoneNumber = item.PhoneNumber
		newUser.FullName = item.FullName
		newUser.Email = item.Email
		newUser.Address = item.Address
		newUser.Gender = item.Gender
		newUser.Birthday = item.Birthday
		newUser.Role = modelUsers.StudentRole

		password, _ := controller.HashedPassword(item.Password)
		newUser.Password = string(password)

		if err := tx.Create(&newUser).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = err.Error()
			return c.JSON(response)
		}
	}

	response.Status = true
	response.Message = config.GetMessageCode("CREATE_SUCCESS")
	return c.JSON(response)
}

// UpdateStudent cập nhật thông tin một Student dựa trên UUID
// @Summary Update student details by UUID
// @Description Update student details by UUID
// @Tags Student
// @Accept json
// @Produce json
// @Param body body []model.UpdateStudent true "Student information to update"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /student [put]
func UpdateStudent(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var payload []*model.UpdateStudent
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := db.Begin()
	defer tx.Commit()

	for _, item := range payload {

		var student model.Student
		if item.UUID != "" {
			if item.IsDeleted {
				// Xóa mềm student bằng cách cập nhật trạng thái
				results := tx.Where("uuid = ?", item.UUID).Delete(&model.Student{})
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("UUID_NOT_FOUND")
					return c.JSON(response)
				}
			} else {
				// Cập nhật student
				results := db.Where("uuid = ?", item.UUID).First(&student)
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("UUID_NOT_FOUND")
					return c.JSON(response)
				}
				updateFields(&student, item)

				if err := tx.Save(&student).Error; err != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("SYSTEM_ERROR")
					return c.JSON(response)
				}
			}
		} else {
			// Tạo student mới
			newUser := new(model.Student)
			newUser.Code = item.Code
			newUser.FirstName = item.FirstName
			newUser.LastName = item.LastName
			newUser.Image = item.Image
			newUser.PhoneNumber = item.PhoneNumber
			newUser.FullName = item.FullName
			newUser.Email = item.Email
			newUser.Address = item.Address
			newUser.Gender = item.Gender
			newUser.Birthday = item.Birthday
			newUser.Role = modelUsers.StudentRole

			password, _ := controller.HashedPassword(item.Password)
			newUser.Password = string(password)

			if err := tx.Create(newUser).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = config.GetMessageCode("SYSTEM_ERROR")
				return c.JSON(response)
			}
		}
	}

	response.Status = true
	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	return c.JSON(response)
}

// DeleteStudent xóa một Student dựa trên UUID
// @Summary Xóa Student
// @Description Xóa một Student dựa trên UUID
// @Tags Student
// @Accept json
// @Produce json
// @Param uuid path string true "UUID của Student"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /student/{uuid} [delete]
func DeleteStudentByUUID(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the UUID from the request parameters
	uuid := c.Params("uuid")

	// Start a database transaction
	tx := database.DB.Begin()
	defer tx.Commit()

	var student model.Student
	result := tx.First(&student, "uuid = ?", uuid)
	if result.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_UUID_EXISTS")
		return c.JSON(response)
	}

	if err := tx.Model(&student).Updates(map[string]interface{}{
		// "deleted_by": model.GetUsername(c),
		"deleted_at": time.Now(),
	}).Error; err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		return c.JSON(response)
	}

	response.Status = true
	response.Message = config.GetMessageCode("DELETE_SUCCESS")
	return c.JSON(response)
}

// RestoreStudent khôi phục một Student đã bị xóa dựa trên UUID
// @Summary Khôi phục Student
// @Description Khôi phục một Student đã bị xóa dựa trên UUID
// @Tags Student
// @Accept json
// @Produce json
// @Param uuid path string true "UUID của Student"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /student/restore/{uuid} [put]
func RestoreStudentByUUID(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the UUID from the request parameters
	uuid := c.Params("uuid")

	// Check if the record exists
	var student model.Student
	results := database.DB.Unscoped().First(&student, "uuid = ?", uuid)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_FOUND")
		return c.JSON(response)
	}

	// Check if the record is already restored
	if !student.Header.DeletedAt.Valid {
		response.Status = false
		response.Message = config.GetMessageCode("ALREADY_RESTORED")
		return c.JSON(response)
	}

	database.DB.Model(&model.Student{}).Unscoped().Where("uuid = ?", uuid).Updates(map[string]interface{}{
		"deleted_at": nil,
		"deleted_by": "",
	})

	response.Status = true
	response.Message = config.GetMessageCode("RESTORE_SUCCESS")
	return c.JSON(response)
}

// CreateTestStudents tạo dữ liệu sinh viên thử nghiệm
// @Summary Create test student data
// @Description Create test student data for development or testing purposes
// @Tags Student
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /student/create-test [post]
func CreateTestStudents(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Mảng JSON chứa dữ liệu sinh viên thử nghiệm
	testStudentsJSON := []byte(`
	[
		
	  ]
	  
	
	`)

	var testStudents []model.Student

	if err := json.Unmarshal(testStudentsJSON, &testStudents); err != nil {
		response.Status = false
		response.Message = "Failed to unmarshal test student data"
		return c.JSON(response)
	}

	tx := database.DB.Begin()
	defer tx.Commit()

	for _, student := range testStudents {
		// Thêm logic kiểm tra trường hợp lỗi hoặc xử lý dữ liệu ở đây (tùy theo yêu cầu)

		// Tạo mới sinh viên thử nghiệm trong cơ sở dữ liệu
		if err := tx.Create(&student).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = "Failed to create test student data"
			return c.JSON(response)
		}
	}

	response.Status = true
	response.Message = "Test student data created successfully"
	return c.JSON(response)
}

func updateFields(student *model.Student, update *model.UpdateStudent) {
	if update.Code != "" {
		student.Code = update.Code
	}
	if update.FirstName != "" {
		student.FirstName = update.FirstName
	}
	if update.LastName != "" {
		student.LastName = update.LastName
	}
	if update.Image != "" {
		student.Image = update.Image
	}
	if update.PhoneNumber != "" {
		student.PhoneNumber = update.PhoneNumber
	}
	if update.FullName != "" {
		student.FullName = update.FullName
	}
	if update.Email != "" {
		student.Email = update.Email
	}
	if update.Address != "" {
		student.Address = update.Address
	}
	if update.Gender != "" {
		student.Gender = update.Gender
	}
	if update.Status != student.Status {
		student.Status = update.Status
	}
	if update.Birthday != "" {
		student.Birthday = update.Birthday
	}
	if update.Password != "" {
		password, _ := controller.HashedPassword(update.Password)
		student.Password = string(password)
	}
}
