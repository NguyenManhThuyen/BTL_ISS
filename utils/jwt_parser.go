package utils

import (
	"app/config"
	"errors"
	//"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type TokenData struct {
	ID  uint
	Code string
	Role      int
	Createdat int64
	Expires   int64
}

func extractToken(c *fiber.Ctx) string {
	bearToken := c.Get("bku-token")
	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 && onlyToken[0] == "Bearer" {
		return onlyToken[1]
	}

	return ""
}

func VerifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)

	if len(tokenString) == 0 {
		msg := config.GetMessageCode("TOKEN_INCORRECT")
		return nil, errors.New(msg)
	}

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(config.Config("JWT_SECRET_KEY")), nil
}

func ExtractTokenData(c *fiber.Ctx) (*TokenData, error) {
	token, err := VerifyToken(c)
	if err != nil {
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// Sử dụng type assertion để lấy giá trị uint từ claims
		id, idOk := claims["id"].(float64)
		if !idOk {
			return nil, errors.New("Could not extract uint ID from token")
		}

		return &TokenData{
			ID:        uint(id),
			Code:      claims["code"].(string),
			Role:      int(claims["role"].(float64)),
			Createdat: int64(claims["iat"].(float64)),
			Expires:   int64(claims["exp"].(float64)),
		}, nil
	}

	return nil, err
}



