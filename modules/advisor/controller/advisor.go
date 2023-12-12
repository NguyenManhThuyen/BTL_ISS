package controller

import (
	"app/config"
	"app/controller"
	"app/database"
	"app/modules/advisor/model"
	"encoding/json"
	"time"

	modelUsers "app/modules/users/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func generateRandomUUID() string {
	return uuid.NewString()
}

// @title Advisor API
// @version 1.0
// @description API for managing advisor data
// @termsOfService http://swagger.io/terms/
// @BasePath /advisor
// @schemes http
// @produce json
// @consumes json

// GetAdvisor trả về danh sách tất cả Advisor
// @Summary Get a list of advisors
// @Description Get a list of all advisors
// @Tags Advisor
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /advisor [get]
func GetAdvisor(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var advisors []model.Advisor
	results := database.DB.Order("ID").Find(&advisors)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = advisors
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetAdvisorByUUID trả về thông tin của một Advisor dựa trên UUID
// @Summary Get advisor by UUID
// @Description Get advisor details by UUID
// @Tags Advisor
// @Produce json
// @Param uuid path string true "Advisor UUID"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /advisor/{uuid} [get]
func GetAdvisorByUUID(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var advisor model.Advisor
	results := database.DB.First(&advisor, "uuid = ?", c.Params("uuid"))
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = advisor
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetAdvisorByMSCB trả về thông tin của một Advisor dựa trên CODE
// @Summary Get advisor by CODE
// @Description Get advisor details by CODE
// @Tags Advisor
// @Produce json
// @Param code path string true "Advisor CODE"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /advisor/code/{code} [get]
func GetAdvisorByMSCB(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var advisor model.Advisor

	// Lấy CODE từ request parameters
	code := c.Params("code")

	results := database.DB.First(&advisor, "code = ?", code)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}

	response.Data = advisor
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// CreateAdvisor tạo một Advisor mới
// @Summary Create a new advisor
// @Description Create a new advisor
// @Tags Advisor
// @Accept json
// @Produce json
// @Param body body []model.CreateAdvisor true "Advisor information"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /advisor [post]
func CreateAdvisor(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	var payload []*model.CreateAdvisor
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := database.DB.Begin()
	defer tx.Commit()

	for _, item := range payload {
		// Thêm logic kiểm tra và xử lý dữ liệu ở đây (tùy thuộc vào yêu cầu)
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
			response.Status = false
			response.Message = config.GetMessageCode("MISSING_FIELDS")
			response.ValidateError = missingFields
			return c.JSON(response)
		}

		// Tạo mới Advisor trong cơ sở dữ liệu
		newUser := new(model.Advisor)
		newUser.ID = uint(item.ID)
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
		newUser.Role = modelUsers.AdvisorRole
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

// UpdateAdvisor cập nhật thông tin một Advisor dựa trên UUID
// @Summary Update advisor details by UUID
// @Description Update advisor details by UUID
// @Tags Advisor
// @Accept json
// @Produce json
// @Param body body []model.UpdateAdvisor true "Advisor information to update"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /advisor [put]
func UpdateAdvisor(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var payload []*model.UpdateAdvisor
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := db.Begin()
	defer tx.Commit()

	for _, item := range payload {
		// Thêm logic kiểm tra và xử lý dữ liệu ở đây (tùy thuộc vào yêu cầu)

		var advisor model.Advisor
		if item.UUID != "" {
			if item.IsDeleted {
				// Xóa mềm Advisor bằng cách cập nhật trạng thái
				results := tx.Where("uuid = ?", item.UUID).Delete(&model.Advisor{})
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("UUID_NOT_FOUND")
					return c.JSON(response)
				}
			} else {
				results := tx.Where("uuid = ?", item.UUID).First(&advisor)
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("UUID_NOT_FOUND")
					return c.JSON(response)
				}

				updateFields(&advisor, item)

				if err := tx.Save(&advisor).Error; err != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("SYSTEM_ERROR")
					return c.JSON(response)
				}
			}
		} else {
			// Tạo mới Advisor
			newUser := new(model.Advisor)
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
			newUser.Role = modelUsers.AdvisorRole

			password, _ := controller.HashedPassword(item.Password)
			newUser.Password = string(password)

			if err := tx.Create(&advisor).Error; err != nil {
				tx.Rollback()
				response.Status = false
				response.Message = err.Error()
				return c.JSON(response)
			}
		}
	}

	response.Status = true
	response.Message = config.GetMessageCode("UPDATE_SUCCESS")
	return c.JSON(response)
}

// Xóa Advisor dựa trên UUID
// @Summary Xóa Advisor
// @Description Xóa một Advisor dựa trên UUID
// @Tags Advisor
// @Accept json
// @Produce json
// @Param uuid path string true "UUID của Advisor"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /advisor/{uuid} [delete]
func DeleteAdvisor(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the ID from the request parameters
	id := c.Params("uuid")

	// Start a database transaction
	tx := database.DB.Begin()
	defer tx.Commit()

	var advisor model.Advisor
	result := tx.First(&advisor, "uuid = ?", id)
	if result.Error != nil {
		tx.Rollback()
		response.Status = false
		response.Message = config.GetMessageCode("NOT_UUID_EXISTS")
		return c.JSON(response)
	}

	if err := tx.Model(&advisor).Updates(map[string]interface{}{
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

// Khôi phục Advisor đã bị xóa
// @Summary Khôi phục Advisor
// @Description Khôi phục một Advisor đã bị xóa
// @Tags Advisor
// @Accept json
// @Produce json
// @Param uuid path string true "UUID của Advisor"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /advisor/restore/{uuid} [put]
func RestoreAdvisor(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Get the UUID from the request parameters
	uuid := c.Params("uuid")

	// Check if the record exists
	var advisor model.Advisor

	results := database.DB.Unscoped().First(&advisor, "uuid = ?", uuid)

	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_FOUND")
		return c.JSON(response)
	}

	// Check if the record is already restored
	if !advisor.Header.DeletedAt.Valid {
		response.Status = false
		response.Message = config.GetMessageCode("ALREADY_RESTORED")
		return c.JSON(response)
	}

	database.DB.Model(&model.Advisor{}).Unscoped().Where("uuid = ?", uuid).Updates(map[string]interface{}{
		"deleted_at": nil,
		"deleted_by": "",
	})

	response.Status = true
	response.Message = config.GetMessageCode("RESTORE_SUCCESS")
	return c.JSON(response)
}

// Tạo dữ liệu Advisor thử nghiệm
// @Summary Create test advisor data
// @Description Create test advisor data for development or testing purposes
// @Tags Advisor
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /advisor/create-test [post]
func CreateTestAdvisors(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// Mảng JSON chứa dữ liệu Advisor thử nghiệm
	testAdvisorsJSON := []byte(`
	[
		
	]
	`)

	var testAdvisors []model.Advisor

	if err := json.Unmarshal(testAdvisorsJSON, &testAdvisors); err != nil {
		response.Status = false
		response.Message = "Failed to unmarshal test advisor data"
		return c.JSON(response)
	}

	tx := database.DB.Begin()
	defer tx.Commit()

	for _, advisor := range testAdvisors {
		// Thêm logic kiểm tra trường hợp lỗi hoặc xử lý dữ liệu ở đây (tùy theo yêu cầu)

		// Tạo mới Advisor thử nghiệm trong cơ sở dữ liệu
		if err := tx.Create(&advisor).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = "Failed to create test advisor data"
			return c.JSON(response)
		}
	}

	response.Status = true
	response.Message = "Test advisor data created successfully"
	return c.JSON(response)
}

func updateFields(advisor *model.Advisor, update *model.UpdateAdvisor) {
	if update.Code != "" {
		advisor.Code = update.Code
	}
	if update.FirstName != "" {
		advisor.FirstName = update.FirstName
	}
	if update.LastName != "" {
		advisor.LastName = update.LastName
	}
	if update.Image != "" {
		advisor.Image = update.Image
	}
	if update.PhoneNumber != "" {
		advisor.PhoneNumber = update.PhoneNumber
	}
	if update.FullName != "" {
		advisor.FullName = update.FullName
	}
	if update.Email != "" {
		advisor.Email = update.Email
	}
	if update.Address != "" {
		advisor.Address = update.Address
	}
	if update.Gender != "" {
		advisor.Gender = update.Gender
	}
	if update.Birthday != "" {
		advisor.Birthday = update.Birthday
	}
	if update.Password != "" {
		password, _ := controller.HashedPassword(update.Password)
		advisor.Password = string(password)
	}
}
