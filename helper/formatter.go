package helper

import (
	"bwastartup/users"
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

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

func FormatValidationError(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}

func GetSign(types string, player string) [16]byte {

	YGG_KEY := "gY7ychHF3AjAKc2u4Fy"
	YGG_TOP_ORG := "gY7ychHF3AjAKc2u4Fy"
	YGG_ORG := "gY7ychHF3AjAKc2u4Fy"

	signPlayer := player + "" + YGG_KEY
	signMerchant := player + "" + YGG_TOP_ORG + "" + YGG_ORG + "" + YGG_KEY

	hashPlayer := md5.Sum([]byte(signPlayer))
	hashMerchant := md5.Sum([]byte(signMerchant))

	typ := hashMerchant
	if types == "player" {
		typ = hashPlayer
	}

	return typ
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
