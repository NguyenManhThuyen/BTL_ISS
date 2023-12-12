package controller

import (
	"app/config"
	"app/controller"
	"app/database"
	"app/modules/facultyOffice/model"
	"encoding/json"
	"time"

	modelUsers "app/modules/users/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func generateRandomUUID() string {
	return uuid.NewString()
}

// GetFacultyOffices trả về danh sách tất cả FacultyOffices
// @Summary Get all FacultyOffices
// @Description Get a list of all FacultyOffices
// @Tags FacultyOffice
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /facultyoffice [get]
func GetFacultyOffices(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var facultyOffices []model.FacultyOffice
	results := database.DB.Select("*").Order("id").Find(&facultyOffices)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = facultyOffices
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetFacultyOfficeByUUID trả về thông tin của một FacultyOffice dựa trên UUID
// @Summary Get a FacultyOffice by UUID
// @Description Get the information of a FacultyOffice by UUID
// @Tags FacultyOffice
// @Param uuid path int true "FacultyOffice UUID"
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /facultyoffice/{uuid} [get]
func GetFacultyOfficeByUUID(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var facultyOffice model.FacultyOffice
	results := database.DB.First(&facultyOffice, "uuid = ?", c.Params("uuid"))
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = facultyOffice
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetFacultyOfficeByMSCB trả về thông tin của một FacultyOffice dựa trên CODE
// @Summary Get FacultyOffice by CODE
// @Description Get the information of a FacultyOffice by CODE
// @Tags FacultyOffice
// @Produce json
// @Param code path string true "Faculty Office CODE"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /facultyoffice/code/{code} [get]
func GetFacultyOfficeByMSCB(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	code := c.Params("code") // Lấy giá trị CODE từ tham số trong URL

	var facultyOffice model.FacultyOffice
	results := database.DB.Where("code = ?", code).First(&facultyOffice)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}

	response.Data = facultyOffice
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// CreateFacultyOffice tạo một FacultyOffice mới
// @Summary Create a new FacultyOffice
// @Description Create a new FacultyOffice
// @Tags FacultyOffice
// @Accept json
// @Produce json
// @Param body body []model.CreateFacultyOffice true "FacultyOffice information"
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /facultyoffice [post]
func CreateFacultyOffice(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	var payload []*model.CreateFacultyOffice
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

		// Tạo FacultyOffice từ dữ liệu được gửi lên
		newUser := new(model.FacultyOffice)
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
		newUser.Role = modelUsers.FacultyOfficeRole

		password, _ := controller.HashedPassword(item.Password)
		newUser.Password = string(password)

		newUser.Header.CreatedAt = time.Now()

		// Thực hiện tạo mới FacultyOffice trong database
		if err := tx.Create(&newUser).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = config.GetMessageCode("CREATE_FAIL")
			return c.JSON(response)
		}
	}

	response.Status = true
	response.Message = config.GetMessageCode("CREATE_SUCCESS")
	return c.JSON(response)
}

// UpdateFacultyOffice cập nhật thông tin một FacultyOffice dựa trên UUID
// @Summary Update a FacultyOffice by UUID
// @Description Update the information of a FacultyOffice by UUID
// @Tags FacultyOffice
// @Accept json
// @Produce json
// @Param body body []model.UpdateFacultyOffice true "FacultyOffice information to update"
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /facultyoffice [put]
func UpdateFacultyOffice(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var payload []*model.UpdateFacultyOffice
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := db.Begin()
	defer tx.Commit()

	for _, item := range payload {

		var facultyOffice model.FacultyOffice
		// Kiểm tra nếu có ID, thực hiện cập nhật
		if item.UUID != "" {
			if item.IsDeleted {
				// Xóa mềm FacultyOffice bằng cách cập nhật trạng thái
				result := tx.Where("uuid = ?", item.UUID).Delete(&model.FacultyOffice{})
				if result.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("UUID_NOT_FOUND")
					return c.JSON(response)
				}
			} else {
				// Cập nhật FacultyOffice
				results := db.Where("uuid = ?", item.UUID).First(&facultyOffice)
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("UUID_NOT_FOUND")
					return c.JSON(response)
				}
				facultyOffice.Code = item.Code
				facultyOffice.FirstName = item.FirstName
				facultyOffice.LastName = item.LastName
				facultyOffice.Image = item.Image
				facultyOffice.PhoneNumber = item.PhoneNumber
				facultyOffice.FullName = item.FullName
				facultyOffice.Email = item.Email
				facultyOffice.Address = item.Address
				facultyOffice.Gender = item.Gender
				facultyOffice.Birthday = item.Birthday
				facultyOffice.Role = modelUsers.FacultyOfficeRole

				password, _ := controller.HashedPassword(item.Password)
				facultyOffice.Password = string(password)

				facultyOffice.Header.UpdatedAt = time.Now()

				if err := tx.Save(&facultyOffice).Error; err != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("SYSTEM_ERROR")
					return c.JSON(response)
				}
			}
			// Tìm FacultyOffice trong database
			results := tx.First(&facultyOffice, item.UUID)
			if results.Error != nil {
				tx.Rollback()
				response.Status = false
				response.Message = config.GetMessageCode("NOT_ID_EXISTS")
				return c.JSON(response)
			}
		} else {
			// Tạo FacultyOffice từ dữ liệu được gửi lên
			newUser := new(model.FacultyOffice)
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
			newUser.Role = modelUsers.FacultyOfficeRole

			password, _ := controller.HashedPassword(item.Password)
			newUser.Password = string(password)

			newUser.Header.CreatedAt = time.Now()
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

// DeleteFacultyOffice xóa một FacultyOffice dựa trên UUID
// @Summary Delete a FacultyOffice by UUID
// @Description Delete a FacultyOffice by UUID
// @Tags FacultyOffice
// @Accept json
// @Param uuid path string true "UUID của FacultyOffice"
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /facultyoffice/{uuid} [delete]
func DeleteFacultyOfficeUUID(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Lấy UUID từ tham số yêu cầu
	uuid := c.Params("uuid")

	// Bắt đầu giao dịch cơ sở dữ liệu
	tx := database.DB.Begin()
	defer tx.Commit()

	var facultyOffice model.FacultyOffice
	result := tx.First(&facultyOffice, "uuid = ?", uuid)
	if result.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_ID_EXISTS")
		return c.JSON(response)
	}

	// Cập nhật trường deleted_at để đánh dấu bản ghi đã bị xóa
	if err := tx.Model(&facultyOffice).Updates(map[string]interface{}{
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

// RestoreFacultyOffice khôi phục một FacultyOffice đã bị xóa dựa trên UUID
// @Summary Restore a deleted FacultyOffice
// @Description Restore a deleted FacultyOffice
// @Param uuid path string true "UUID của FacultyOffice"
// @Tags FacultyOffice
// @Accept json
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /facultyoffice/restore/{uuid} [put]
func RestoreFacultyOfficeUUID(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Lấy ID từ tham số yêu cầu
	uuid := c.Params("uuid")

	// Kiểm tra xem bản ghi có tồn tại không
	var facultyOffice model.FacultyOffice
	results := database.DB.Unscoped().First(&facultyOffice, "uuid = ?", uuid)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_FOUND")
		return c.JSON(response)
	}

	// Kiểm tra xem bản ghi đã được khôi phục chưa
	if !facultyOffice.Info.DeletedAt.Valid {
		response.Status = false
		response.Message = config.GetMessageCode("ALREADY_RESTORED")
		return c.JSON(response)
	}

	// Sử dụng Unscoped để khôi phục bản ghi đã bị xóa
	database.DB.Model(&model.FacultyOffice{}).Unscoped().Where("uuid = ?", uuid).Updates(map[string]interface{}{
		"deleted_at": nil,
	})

	response.Status = true
	response.Message = config.GetMessageCode("RESTORE_SUCCESS")
	return c.JSON(response)
}

// CreateTestFacultyOffices creates test FacultyOffice data for development or testing purposes.
// @Summary Create test FacultyOffices
// @Description Create test FacultyOffices
// @Tags FacultyOffice
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /facultyoffice/create-test [post]
func CreateTestFacultyOffices(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// JSON array containing test FacultyOffice data
	testFacultyOfficesJSON := []byte(`
	[
		
	  ]
	  
    `)

	var testFacultyOffices []*model.CreateFacultyOffice
	if err := json.Unmarshal(testFacultyOfficesJSON, &testFacultyOffices); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		return c.JSON(response)
	}

	tx := database.DB.Begin()
	defer tx.Commit()

	for _, item := range testFacultyOffices {
		// Tạo FacultyOffice từ dữ liệu được gửi lên
		newFacultyOffice := new(model.FacultyOffice)

		newFacultyOffice.FirstName = item.FirstName
		newFacultyOffice.LastName = item.LastName
		newFacultyOffice.Image = item.Image
		newFacultyOffice.LastName = item.LastName
		newFacultyOffice.PhoneNumber = item.PhoneNumber
		newFacultyOffice.FullName = item.FullName
		newFacultyOffice.Email = item.Email
		newFacultyOffice.Address = item.Address
		newFacultyOffice.Gender = item.Gender
		newFacultyOffice.Birthday = item.Birthday
		newFacultyOffice.Role = 4
		newFacultyOffice.Info.CreatedAt = time.Now()
		if err := tx.Create(&newFacultyOffice).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = config.GetMessageCode("SYSTEM_ERROR")
			return c.JSON(response)
		}
	}

	response.Status = true
	response.Message = config.GetMessageCode("CREATE_SUCCESS")
	return c.JSON(response)
}
