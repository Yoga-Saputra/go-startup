package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"startup/app/users"
	"startup/config"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// helper user
type UserFormatter struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
}

func FormatUser(user users.User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		Token:      token,
	}
	return formatter
}

func FormatValidationError(err error) interface{} {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	errorMessage := gin.H{"errors": errors}
	return errorMessage
}

func ErrorValidation(err error, c *gin.Context, msg string, status string, code int, errMsg interface{}) {
	jsonConvert := InterfaceToJson(errMsg)
	log := ApiResponse(msg, http.StatusUnprocessableEntity, status, string(jsonConvert))
	config.Loggers("error", log)

	response := ApiResponse(msg, http.StatusUnprocessableEntity, status, errMsg)
	c.JSON(code, response)
}

// Convert Map to JSON
func MapToJson(params map[string]string) []byte {
	jsonContent, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return jsonContent
	}
	return jsonContent
}

func InterfaceToJson(params interface{}) []byte {
	jsonContent, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return jsonContent
	}
	return jsonContent
}

func InterfaceToInt(byteCode interface{}) (int, error) {
	code, _ := json.Marshal(byteCode)
	// convert []byte to int
	codeInt, err := strconv.Atoi(string(code))
	if err != nil {
		return codeInt, err
	}

	return codeInt, nil
}

func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert interface to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map

	if err != nil {
		return newMap, err
	}

	return
}
