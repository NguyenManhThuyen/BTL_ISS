package controller

import "golang.org/x/crypto/bcrypt"

// func Test(c *fiber.Ctx) string {
// 	response := new(config.DataResponse)

// 	tokenData, err := utils.ExtractTokenData(c)
// 	if err != nil {
// 		response.Status = false
// 		response.Message = config.GetMessageCode("ERROR_GET_USERNAME")
// 	}

// 	return tokenData.Username
// }

func HashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, err
}