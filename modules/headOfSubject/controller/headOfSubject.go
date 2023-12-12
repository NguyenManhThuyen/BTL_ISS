package controller

import (
	"app/config"
	"app/controller"
	"app/database"
	"app/modules/headOfSubject/model"
	"encoding/json"
	"time"

	modelUsers "app/modules/users/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func generateRandomUUID() string {
	return uuid.NewString()
}

// GetHeadOfSubject trả về danh sách tất cả HeadOfSubject
// @Summary Get all HeadOfSubjects
// @Description Get a list of all HeadOfSubjects
// @Tags HeadOfSubject
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /headofsubject [get]
func GetHeadOfSubject(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var headsOfSubject []model.HeadOfSubject
	results := database.DB.Select("*").Order("id").Find(&headsOfSubject)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = headsOfSubject
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetHeadOfSubjectByUUID trả về thông tin của một HeadOfSubject dựa trên UUID
// @Summary Get a HeadOfSubject by UUID
// @Description Get the information of a HeadOfSubject by UUID
// @Tags HeadOfSubject
// @Param id path int true "HeadOfSubject UUID"
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /headofsubject/{uuid} [get]
func GetHeadOfSubjectByUUID(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var headOfSubject model.HeadOfSubject
	results := database.DB.First(&headOfSubject, "uuid = ?", c.Params("uuid"))
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = headOfSubject
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetHeadOfSubjectByMSCB trả về thông tin của một HeadOfSubject dựa trên CODE
// @Summary Get HeadOfSubject by CODE
// @Description Get the information of a HeadOfSubject by CODE
// @Tags HeadOfSubject
// @Produce json
// @Param code path string true "Head Of Subject CODE"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /headofsubject/code/{code} [get]
func GetHeadOfSubjectByMSCB(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	code := c.Params("code") // Lấy giá trị CODE từ tham số trong URL

	var headOfSubject model.HeadOfSubject
	results := database.DB.Where("code = ?", code).First(&headOfSubject)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}

	response.Data = headOfSubject
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// CreateHeadOfSubject tạo một HeadOfSubject mới
// @Summary Create a new HeadOfSubject
// @Description Create a new HeadOfSubject
// @Tags HeadOfSubject
// @Accept json
// @Produce json
// @Param body body []model.CreateHeadOfSubject true "HeadOfSubject information"
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /headofsubject [post]
func CreateHeadOfSubject(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	var payload []*model.CreateHeadOfSubject
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

		newUser := new(model.HeadOfSubject)
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
		newUser.Role = modelUsers.HeadOfSubjectRole

		password, _ := controller.HashedPassword(item.Password)
		newUser.Password = string(password)

		newUser.Header.CreatedAt = time.Now()
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

// UpdateHeadOfSubject cập nhật thông tin một HeadOfSubject dựa trên UUID
// @Summary Update a HeadOfSubject by UUID
// @Description Update the information of a HeadOfSubject by UUID
// @Tags HeadOfSubject
// @Accept json
// @Produce json
// @Param body body []model.UpdateHeadOfSubject true "HeadOfSubject information to update"
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /headofsubject [put]
func UpdateHeadOfSubject(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var payload []*model.UpdateHeadOfSubject
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := db.Begin()
	defer tx.Commit()

	for _, item := range payload {

		var headOfSubject model.HeadOfSubject

		if item.UUID != "" {
			if item.IsDeleted {
				// Xóa mềm headOfSubject bằng cách cập nhật trạng thái
				results := tx.Where("uuid = ?", item.UUID).Delete(&model.HeadOfSubject{})
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("UUID_NOT_FOUND")
					return c.JSON(response)
				}
			} else {
				// Cập nhật HeadOfSubject
				results := db.Where("uuid = ?", item.UUID).First(&headOfSubject)
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("UUID_NOT_FOUND")
					return c.JSON(response)
				}

				headOfSubject.Code = item.Code
				headOfSubject.FirstName = item.FirstName
				headOfSubject.LastName = item.LastName
				headOfSubject.Image = item.Image
				headOfSubject.PhoneNumber = item.PhoneNumber
				headOfSubject.FullName = item.FullName
				headOfSubject.Email = item.Email
				headOfSubject.Address = item.Address
				headOfSubject.Gender = item.Gender
				headOfSubject.Birthday = item.Birthday
				headOfSubject.Role = modelUsers.HeadOfSubjectRole

				password, _ := controller.HashedPassword(item.Password)
				headOfSubject.Password = string(password)

				headOfSubject.Header.UpdatedAt = time.Now()

				if err := tx.Save(&headOfSubject).Error; err != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("SYSTEM_ERROR")
					return c.JSON(response)
				}
			}
		} else {
			// Tạo mới headOfSubject
			newUser := new(model.HeadOfSubject)
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
			newUser.Role = modelUsers.HeadOfSubjectRole

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

// DeleteHeadOfSubject xóa một HeadOfSubject dựa trên UUID
// @Summary Delete a HeadOfSubject by UUID
// @Description Delete a HeadOfSubject by UUID
// @Tags HeadOfSubject
// @Accept json
// @Param uuid path int true "HeadOfSubject UUID"
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /headofsubject/{uuid} [delete]
func DeleteHeadOfSubjectByUUID(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the ID from the request parameters
	uuid := c.Params("uuid")

	// Start a database transaction
	tx := database.DB.Begin()
	defer tx.Commit()

	var headOfSubject model.HeadOfSubject
	result := tx.First(&headOfSubject, "uuid = ?", uuid)
	if result.Error != nil {
		tx.Rollback()
		response.Status = false
		response.Message = config.GetMessageCode("NOT_ID_EXISTS")
		return c.JSON(response)
	}

	if err := tx.Model(&headOfSubject).Updates(map[string]interface{}{
		// "deleted_by": model.GetUsername(c),
		"deleted_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		response.Status = false
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		return c.JSON(response)
	}

	response.Status = true
	response.Message = config.GetMessageCode("DELETE_SUCCESS")
	return c.JSON(response)
}

// RestoreHeadOfSubject khôi phục một HeadOfSubject đã bị xóa dựa trên UUID
// @Summary Restore a deleted HeadOfSubject
// @Description Khôi phục HeadOfSubject đã bị xóa dựa trên UUID
// @Param id path int true "HeadOfSubject UUID"
// @Tags HeadOfSubject
// @Accept json
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /headofsubject/restore/{uuid} [post]
func RestoreHeadOfSubjectByUUID(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the ID from the request parameters
	uuid := c.Params("uuid")

	// Check if the record exists
	var headOfSubject model.HeadOfSubject
	results := database.DB.Unscoped().First(&headOfSubject, "uuid = ?", uuid)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_FOUND")
		return c.JSON(response)
	}

	// Check if the record is already restored
	if !headOfSubject.Info.DeletedAt.Valid {
		response.Status = false
		response.Message = config.GetMessageCode("ALREADY_RESTORED")
		return c.JSON(response)
	}

	database.DB.Model(&model.HeadOfSubject{}).Unscoped().Where("uuid = ?", uuid).Updates(map[string]interface{}{
		"deleted_at": nil,
		"deleted_by": "",
	})

	response.Status = true
	response.Message = config.GetMessageCode("RESTORE_SUCCESS")
	return c.JSON(response)
}

// CreateTestHeadOfSubjects tạo dữ liệu HeadOfSubject thử nghiệm
// @Summary Create test HeadOfSubjects
// @Description Create test HeadOfSubjects
// @Tags HeadOfSubject
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /headofsubject/create-test [post]
func CreateTestHeadOfSubjects(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Mảng JSON chứa dữ liệu HeadOfSubject thử nghiệm
	testHeadOfSubjectsJSON := []byte(`
	[
		
	  ]
	  
`)

	var testHeadOfSubjects []*model.CreateHeadOfSubject
	if err := json.Unmarshal(testHeadOfSubjectsJSON, &testHeadOfSubjects); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("SYSTEM_ERROR")
		return c.JSON(response)
	}

	tx := database.DB.Begin()
	defer tx.Commit()

	for _, item := range testHeadOfSubjects {
		newUser := new(model.HeadOfSubject)
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
		newUser.Role = modelUsers.HeadOfSubjectRole

		password, _ := controller.HashedPassword(item.Password)
		newUser.Password = string(password)

		newUser.Header.CreatedAt = time.Now()

		if err := tx.Create(&newUser).Error; err != nil {
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
