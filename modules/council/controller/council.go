package controller

import (
	"app/config"
	"app/controller"
	"app/database"
	"app/modules/council/model"
	"encoding/json"
	"time"

	modelUsers "app/modules/users/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func generateRandomUUID() string {
	return uuid.NewString()
}

// @title Council API
// @version 1.0
// @description API for managing council data
// @termsOfService http://swagger.io/terms/
// @BasePath /council
// @schemes http
// @produce json
// @consumes json

// GetCouncils get a list of all councils.
// @Summary Get a list of councils
// @Description Get a list of councils.
// @Tags Council
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /council [get]
func GetCouncils(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var councils []model.Council
	results := database.DB.Select("*").Order("id").Find(&councils)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = councils
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetCouncilByUUID returns council details by UUID.
// @Summary Get council by UUID
// @Description Get council details by UUID.
// @Tags Council
// @Produce json
// @Param uuid path string true "Council UUID"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /council/{uuid} [get]
func GetCouncilByUUID(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var council model.Council
	results := database.DB.First(&council, "uuid = ?", c.Params("uuid"))
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = council
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// GetCouncilByMSCB trả về thông tin của một CBHD dựa trên CODE
// @Summary Get council by CODE
// @Description Get council details by CODE
// @Tags Council
// @Produce json
// @Param code path string true "Council CODE"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /council/code/{code} [get]
func GetCouncilByMSCB(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	var council model.Council
	results := database.DB.First(&council, "code = ?", c.Params("code"))
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("GET_DATA_FAIL")
		return c.JSON(response)
	}
	response.Data = council
	response.Status = true
	response.Message = config.GetMessageCode("GET_DATA_SUCCESS")
	return c.JSON(response)
}

// CreateCouncil creates a new council.
// @Summary Create a new council
// @Description Create a new council.
// @Tags Council
// @Accept json
// @Produce json
// @Param body body []model.CreateCouncil true "Council information"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /council [post]
func CreateCouncil(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	var payload []*model.CreateCouncil
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

		newUser := new(model.Council)
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
		newUser.Role = modelUsers.CouncilRole

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

// UpdateCouncil updates council details by UUID.
// @Summary Update council details by UUID
// @Description Update council details by UUID.
// @Tags Council
// @Accept json
// @Produce json
// @Param body body []model.UpdateCouncil true "Council information to update"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /council [put]
func UpdateCouncil(c *fiber.Ctx) error {
	response := new(config.DataResponse)
	db := database.DB

	var payload []*model.UpdateCouncil
	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := db.Begin()
	defer tx.Commit()

	for _, item := range payload {
		// You can implement the logic to update the Council model based on the UpdateCouncil struct here
		listCheck := []string{"Code", "FirstName", "LastName", "FullName", "Email", "PhoneNumber", "Birthday"}
		vItem := map[string]string{
			"Code":        item.Code,
			"FirstName":   item.FirstName,
			"LastName":    item.LastName,
			"FullName":    item.FullName,
			"Email":       item.Email,
			"PhoneNumber": item.PhoneNumber,
			"Birthday":    item.Birthday,
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

		var council model.Council
		if item.UUID != "" {
			if item.IsDeleted {
				// Xóa mềm council bằng cách cập nhật trạng thái
				results := tx.Where("uuid = ?", item.UUID).Delete(&model.Council{})
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("UUID_NOT_FOUND")
					return c.JSON(response)
				}
			} else {
				results := tx.Where("uuid = ?", item.UUID).First(&council)
				if results.Error != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("UUID_NOT_FOUND")
					return c.JSON(response)
				}

				council.Code = item.Code
				council.FirstName = item.FirstName
				council.LastName = item.LastName
				council.Image = item.Image
				council.PhoneNumber = item.PhoneNumber
				council.FullName = item.FullName
				council.Email = item.Email
				council.Address = item.Address
				council.Gender = item.Gender
				council.Birthday = item.Birthday
				council.Role = modelUsers.CouncilRole

				password, _ := controller.HashedPassword(item.Password)
				council.Password = string(password)

				council.Header.UpdatedAt = time.Now()

				if err := tx.Save(&council).Error; err != nil {
					tx.Rollback()
					response.Status = false
					response.Message = config.GetMessageCode("SYSTEM_ERROR")
					return c.JSON(response)
				}
			}
		} else {

			newUser := new(model.Council)
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
			newUser.Role = modelUsers.CouncilRole

			password, _ := controller.HashedPassword(item.Password)
			newUser.Password = string(password)

			newUser.Header.CreatedAt = time.Now()

			// Call Save function to automatically invoke BeforeUpdate before updating
			if err := tx.Save(&newUser).Error; err != nil {
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

// DeleteCouncil deletes a council by UUID.
// @Summary Delete Council
// @Description Delete a council by UUID.
// @Tags Council
// @Accept json
// @Produce json
// @Param uuid path string true "UUID of Council"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /council/{uuid} [delete]
func DeleteCouncil(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	uuid := c.Params("uuid")

	tx := database.DB.Begin()
	defer tx.Commit()

	var council model.Council
	result := tx.First(&council, "uuid = ?", uuid)
	if result.Error != nil {
		tx.Rollback()
		response.Status = false
		response.Message = config.GetMessageCode("NOT_ID_EXISTS")
		return c.JSON(response)
	}

	if err := tx.Model(&council).Updates(map[string]interface{}{
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

// RestoreCouncil restores a deleted council by UUID.
// @Summary Restore Council
// @Description Restore a deleted council by UUID.
// @Tags Council
// @Accept json
// @Produce json
// @Param uuid path string true "UUID of Council"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /council/restore/{uuid} [put]
func RestoreCouncil(c *fiber.Ctx) error {

	response := new(config.DataResponse)

	uuid := c.Params("uuid")

	var council model.Council
	results := database.DB.Unscoped().First(&council, "uuid = ?", uuid)
	if results.Error != nil {
		response.Status = false
		response.Message = config.GetMessageCode("NOT_FOUND")
		return c.JSON(response)
	}

	if !council.Header.DeletedAt.Valid {
		response.Status = false
		response.Message = config.GetMessageCode("ALREADY_RESTORED")
		return c.JSON(response)
	}

	database.DB.Model(&model.Council{}).Unscoped().Where("uuid = ?", uuid).Updates(map[string]interface{}{
		"deleted_at": nil,
		"deleted_by": "",
	})

	response.Status = true
	response.Message = config.GetMessageCode("RESTORE_SUCCESS")
	return c.JSON(response)

}

// CreateTestCouncil creates test council data for development or testing purposes.
// @Summary Create test council data
// @Description Create test council data for development or testing purposes.
// @Tags Council
// @Produce json
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /council/create-test [post]
func CreateTestCouncil(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	testCouncilsJSON := []byte(`
	[
		
	  ]
	  
    `)

	var testCouncils []model.Council

	if err := json.Unmarshal(testCouncilsJSON, &testCouncils); err != nil {
		response.Status = false
		response.Message = "Failed to unmarshal test council data"
		return c.JSON(response)
	}

	tx := database.DB.Begin()
	defer tx.Commit()

	for _, council := range testCouncils {
		// Thêm logic kiểm tra trường hợp lỗi hoặc xử lý dữ liệu ở đây (tùy theo yêu cầu)

		// Tạo mới sinh viên thử nghiệm trong cơ sở dữ liệu
		if err := tx.Create(&council).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = "Failed to create test council data"
			return c.JSON(response)
		}
	}

	response.Status = true
	response.Message = "Test council data created successfully"
	return c.JSON(response)
}
