package controller

import (
	// "fmt"
	"strings"

	// "time"
	"app/config"
	"app/database"
	"app/utils"

	controllerAdvisor "app/modules/advisor/controller"
	controllerCouncil "app/modules/council/controller"
	controllerFacultyOffice "app/modules/facultyOffice/controller"
	controllerHeadOfSubject "app/modules/headOfSubject/controller"
	controllerStudent "app/modules/student/controller"

	modelAdvisor "app/modules/advisor/model"
	modelCouncil "app/modules/council/model"
	modelFacultyOffice "app/modules/facultyOffice/model"
	modelHeadOfSubject "app/modules/headOfSubject/model"
	modelStudent "app/modules/student/model"
	modelUsers "app/modules/users/model"

	"github.com/gofiber/fiber/v2"

	// "github.com/golang-jwt/jwt"
	// "github.com/wpcodevo/golang-fiber-jwt/initializers"
	// "github.com/wpcodevo/golang-fiber-jwt/models"

)

// SignUpUser đăng ký một người dùng mới.
// @Summary Đăng ký người dùng mới
// @Description Đăng ký người dùng mới với thông tin được cung cấp.
// @Tags User
// @Accept json
// @Produce json
// @Param body body []modelUsers.SignUpInput true "Thông tin đăng ký người dùng"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /signup [post]
// @Security ApiKeyAuth
func SignUpUser(c *fiber.Ctx) error {
	var payload []*modelUsers.SignUpInput

	response := new(config.DataResponse)
	response.Status = false

	if err := c.BodyParser(&payload); err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("PARAM_ERROR")
		return c.JSON(response)
	}

	tx := database.DB.Begin()
	defer tx.Commit()

	for _, item := range payload {
		if item.Password != item.PasswordConfirm {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Passwords do not match"})
		}
		if item.Role == modelUsers.StudentRole {
			return controllerStudent.CreateStudent(c)
		} else if item.Role == modelUsers.AdvisorRole {
			return controllerAdvisor.CreateAdvisor(c)
		} else if item.Role == modelUsers.CouncilRole {
			return controllerCouncil.CreateCouncil(c)
		} else if item.Role == modelUsers.FacultyOfficeRole {
			return controllerFacultyOffice.CreateFacultyOffice(c)
		} else if item.Role == modelUsers.HeadOfSubjectRole {
			return controllerHeadOfSubject.CreateHeadOfSubject(c)
		}
	}

	return c.JSON(response)
}

// SignInUser đăng nhập người dùng.
// @Summary Đăng nhập người dùng và trả về token.
// @Description Đăng nhập với thông tin đăng nhập được cung cấp và trả về token nếu thành công.
// @Tags User
// @Accept json
// @Produce json
// @Param body body modelUsers.SignInInput true "Thông tin đăng nhập"
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /signin [post]
// @Security ApiKeyAuth
func SignInUser(c *fiber.Ctx) error {
	var payload modelUsers.SignInInput

	response := new(config.DataResponse)
	response.Status = false

	tx := database.DB.Begin()
	defer tx.Commit()

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if errors := modelUsers.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// check validate
	listCheck := []string{"email", "password"}
	payloadMap := map[string]string{"email": payload.Email, "password": payload.Password}
	errors := utils.RequireCheck(listCheck, payloadMap, map[string]string{})

	// Check max length
	listCheck = []string{"email:30", "password:20"}
	errors = utils.MaxLengthCheck(listCheck, payloadMap, errors)

	if len(errors) > 0 {
		response.Message = "validate"
		response.ValidateError = errors
		return c.JSON(response)
	}

	var tokenString string

	if payload.Role == modelUsers.StudentRole {
		var student modelStudent.Student

		if err := tx.Model(&modelStudent.Student{}).First(&student, "EMAIL = ?", strings.ToLower(payload.Email)).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = config.GetMessageCode("INVALID_EMAIL")
			return c.JSON(response)
		}

		// err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(payload.Password))
		// if err != nil {
		// 	tx.Rollback()
		// 	response.Status = false
		// 	response.Message = config.GetMessageCode("INVALID_PASSWORD")
		// 	return c.JSON(response)
		// }

		tokenString, _ = utils.GenerateAccessTokenBKU(student.ID, student.FullName, student.Role, student.Code)
		// if err != nil {
		// 	tx.Rollback()
		// 	response.Status = false
		// 	response.Message = config.GetMessageCode("FAILED_TO_GENERATE_TOKEN")
		// 	return c.JSON(response)
		// }

		// Save session
		// store := database.Store
		// timeExpire := config.Config("JWT_EXPIRED_TIME")
		// minutesCount, _ := strconv.Atoi(timeExpire)
		// store.Set(student.UUID, []byte(tokenString), time.Duration(minutesCount)*time.Minute)

		// // Lưu token vào cookie
		// c.Cookie(&fiber.Cookie{
		// 	Name:     "token",
		// 	Value:    tokenString,
		// 	Path:     "/",
		// 	MaxAge:   60 * 60 * 24, // 1 day in seconds
		// 	Secure:   false,        // Set true if using HTTPS
		// 	HTTPOnly: true,
		// })
		response.Message = "LOGIN SUCCESS"
		response.Status = true
		response.Data = createResultData(&student, tokenString)
	} else if payload.Role == modelUsers.AdvisorRole {
		var advisor modelAdvisor.Advisor

		if err := tx.Model(&modelAdvisor.Advisor{}).First(&advisor, "EMAIL = ?", strings.ToLower(payload.Email)).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = config.GetMessageCode("INVALID_EMAIL_PASSWORD")
			return c.JSON(response)
		}


		tokenString, _ = utils.GenerateAccessTokenBKU(advisor.ID, advisor.FullName, advisor.Role, advisor.Code)


		// Save session

		// // Lưu token vào cookie
		// c.Cookie(&fiber.Cookie{
		// 	Name:     "token",
		// 	Value:    tokenString,
		// 	Path:     "/",
		// 	MaxAge:   60 * 60 * 24, // 1 day in seconds
		// 	Secure:   false,        // Set true if using HTTPS
		// 	HTTPOnly: true,
		// })
		response.Message = "LOGIN SUCCESS"
		response.Status = true
		response.Data = createResultData(&advisor, tokenString)
	} else if payload.Role == modelUsers.CouncilRole {
		var council modelCouncil.Council

		if err := tx.Model(&modelCouncil.Council{}).First(&council, "EMAIL = ?", strings.ToLower(payload.Email)).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = config.GetMessageCode("INVALID_EMAIL_PASSWORD")
			return c.JSON(response)
		}


		tokenString, _ = utils.GenerateAccessTokenBKU(council.ID, council.FullName, council.Role, council.Code)

		// // Lưu token vào cookie
		// c.Cookie(&fiber.Cookie{
		// 	Name:     "token",
		// 	Value:    tokenString,
		// 	Path:     "/",
		// 	MaxAge:   60 * 60 * 24, // 1 day in seconds
		// 	Secure:   false,        // Set true if using HTTPS
		// 	HTTPOnly: true,
		// })
		response.Message = "LOGIN SUCCESS"
		response.Status = true
		response.Data = createResultData(&council, tokenString)
	} else if payload.Role == modelUsers.FacultyOfficeRole {
		var facultyOffice modelFacultyOffice.FacultyOffice

		if err := tx.Model(&modelFacultyOffice.FacultyOffice{}).First(&facultyOffice, "EMAIL = ?", strings.ToLower(payload.Email)).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = config.GetMessageCode("INVALID_EMAIL_PASSWORD")
			return c.JSON(response)
		}



		tokenString, _ = utils.GenerateAccessTokenBKU(facultyOffice.ID, facultyOffice.FullName, facultyOffice.Role, facultyOffice.Code)


		// // Lưu token vào cookie
		// c.Cookie(&fiber.Cookie{
		// 	Name:     "token",
		// 	Value:    tokenString,
		// 	Path:     "/",
		// 	MaxAge:   60 * 60 * 24, // 1 day in seconds
		// 	Secure:   false,        // Set true if using HTTPS
		// 	HTTPOnly: true,
		// })
		response.Message = "LOGIN SUCCESS"
		response.Status = true
		response.Data = createResultData(&facultyOffice, tokenString)
	} else if payload.Role == modelUsers.HeadOfSubjectRole {
		var headOfSubject modelHeadOfSubject.HeadOfSubject

		if err := tx.Model(&modelHeadOfSubject.HeadOfSubject{}).First(&headOfSubject, "EMAIL = ?", strings.ToLower(payload.Email)).Error; err != nil {
			tx.Rollback()
			response.Status = false
			response.Message = config.GetMessageCode("INVALID_EMAIL_PASSWORD")
			return c.JSON(response)
		}



		tokenString, _ = utils.GenerateAccessTokenBKU(headOfSubject.ID, headOfSubject.FullName, headOfSubject.Role, headOfSubject.Code)

		// // Lưu token vào cookie
		// c.Cookie(&fiber.Cookie{
		// 	Name:     "token",
		// 	Value:    tokenString,
		// 	Path:     "/",
		// 	MaxAge:   60 * 60 * 24, // 1 day in seconds
		// 	Secure:   false,        // Set true if using HTTPS
		// 	HTTPOnly: true,
		// })
		response.Message = "LOGIN SUCCESS"
		response.Status = true
		response.Data = createResultData(&headOfSubject, tokenString)
	}

	return c.JSON(response)
}

// LogoutUser đăng xuất người dùng.
// @Summary Đăng xuất người dùng và hủy token.
// @Description Đăng xuất người dùng và hủy token nếu tồn tại.
// @Tags User
// @Success 200 {object} config.DataResponse
// @Failure 500 {object} config.DataResponse
// @Router /logout [post]
func LogoutUser(c *fiber.Ctx) error {
	response := new(config.DataResponse)

	// expired := time.Now().Add(-time.Hour * 24)
	// c.Cookie(&fiber.Cookie{
	// 	Name:    "token",
	// 	Value:   "",
	// 	Expires: expired,
	// })

	// store := database.Store
	// // Lấy giá trị cookie token từ người dùng
	// //tokenValue := c.Cookies("token")

	// // Ép kiểu giá trị từ slice of bytes sang string
	// //email, _ := store.Get(tokenValue)
	// //emailString := string(tokenValue)
	// email := utils.GetTokenUUID(c)

	// // Xóa thông tin người dùng khỏi session
	// err := store.Delete(email)
	// if err != nil {
	// 	response.Status = false
	// 	response.Message = config.GetMessageCode("ERROR_GET_USERNAME")
	// }
	response.Status = true
	response.Message = config.GetMessageCode("LOGOUT_SUCCESS")
	return c.JSON(response)
}

func createResultData(user interface{}, token string) map[string]interface{} {
	resultData := make(map[string]interface{})

	switch user := user.(type) {
	case *modelStudent.Student:
		filteredUser := modelUsers.FilterUserRecordStudent(user)
		resultData["user"] = filteredUser
	case *modelAdvisor.Advisor:
		filteredUser := modelUsers.FilterUserRecordAdvisor(user)
		resultData["user"] = filteredUser
	case *modelCouncil.Council:
		filteredUser := modelUsers.FilterUserRecordCouncil(user)
		resultData["user"] = filteredUser
	case *modelFacultyOffice.FacultyOffice:
		filteredUser := modelUsers.FilterUserRecordFacultyOffice(user)
		resultData["user"] = filteredUser
	case *modelHeadOfSubject.HeadOfSubject:
		filteredUser := modelUsers.FilterUserRecordHeadOfSubject(user)
		resultData["user"] = filteredUser
	}
	resultData["token"] = token

	return resultData
}
