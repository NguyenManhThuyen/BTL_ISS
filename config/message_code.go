package config

var messageList = map[string]string{
	"PARAM_ERROR":                 "MSG_V0000",  // param error
	"MAX_LENGTH":                  "MSG_V0001",  // param max length
	"FIX_LENGTH":                  "MSG_V0002",  // param fix length
	"FORMAT_NUMBER":               "MSG_V0004",  // param is number
	"FORMAT_DATE":                 "MSG_V0003",  // param format date is YYYY-MM-DD. Ex: 2023-01-01
	"REQUIRE":                     "MSG_V0005",  // Param require
	"KEY_NOT_FOUND":               "MSG_S0000",  // key error not found
	"SYSTEM_ERROR":                "MSG_S0001",  // system error
	"TOKEN_INCORRECT":             "MSG_S0002",  // token invalid
	"GET_DATA_FAIL":               "MSG_RE0001", // get data fail
	"CREATE_SUCCESS":              "MSG_CI0001", // Create new data success
	"NOT_ID_EXISTS":               "MSG_RE0002", // No item with that Id exists
	"GET_DATA_SUCCESS":            "MSG_RI0001", // Get data success
	"EMAIL_PASSWORD_INCORRECT": "MSG_N0000",  // User name incorrect
	"MISSING_FIELDS":              "MSG_V1000",  // Missing fields
	"UPDATE_SUCCESS":              "MSG_UI0001", // Update data success
	"DELETE_SUCCESS":              "MSG_DI0001", // Delete data success
	"ALREADY_RESTORED":            "MSG_AR0001", // Already restored
	"INVALID_EMAIL_PASSWORD":      "MSG_N0001",  // INVALID_EMAIL_PASSWORD
	"FAILED_TO_GENERATE_TOKEN":    "MSG_S0003",  // Failed to generate JWT token
	"LOGIN_SUCCESS":               "MSG_S0004",  // LOGIN SUCCESS
	"ERROR_GET_EMAIL":          "MSG_S0005",  // ERROR GET EMAIL
	"LOGOUT_SUCCESS":              "MSG_S0006",  // LOGOUT SUCCESS
	"SIGN_UP_SUCCESS":             "MSG_S0007",  // SIGN UP SUCCESS
	//"RESTORE_SUCCESS":
	//"CONFLICT_EMAIL"
	//"ERROR_END_TIME_ADMIN"
}

func GetMessageCode(key string) string {
	var message = key
	if msg, ok := messageList[key]; ok {
		message = msg
	}

	return message
}
