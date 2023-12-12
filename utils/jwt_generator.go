package utils

import (
	"app/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateAccessToken(userAgent, ipAddress string) (string, error) {
	secret := config.Config("JWT_SECRET_KEY")
	timeExpire := config.Config("JWT_EXPIRED_TIME")

	minutesCount, _ := strconv.Atoi(timeExpire)

	claims := jwt.MapClaims{}

	claims["useragent"] = userAgent
	claims["ipaddress"] = ipAddress
	claims["createdat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func GenerateAccessTokenBKU(id uint,fullName string ,role int, code string) (string, error) {
	secret := config.Config("JWT_SECRET_KEY")
	timeExpire := config.Config("JWT_EXPIRED_TIME")

	minutesCount, _ := strconv.Atoi(timeExpire)

	claims := jwt.MapClaims{}

	claims["id"] = id
	claims["code"] = code
	claims["role"] = role
	//claims["ipaddress"] = ipAddress
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func EncodeDataTokenMobile(employeeId string, dateInSeconds int64, coordinates string, shiftId int) (string, error) {
	secret := config.Config("JWT_DATA_SECRET_KEY")
	timeExpire := config.Config("JWT_DATA_EXPIRED_TIME")

	minutesCount, _ := strconv.Atoi(timeExpire)

	claims := jwt.MapClaims{}

	claims["employeeId"] = employeeId
	claims["dateInSeconds"] = dateInSeconds
	claims["coordinates"] = coordinates
	claims["shiftId"] = shiftId
	claims["createdat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
