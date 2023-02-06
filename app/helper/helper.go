package helper

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	response struct {
		Meta meta        `json:"meta"`
		Data interface{} `json:"data"`
	}
	meta struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Status  string `json:"status"`
	}
)

func ApiResponse(message string, code int, status string, data interface{}) interface{} {
	meta := meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	response := response{
		Meta: meta,
		Data: data,
	}

	// data, err := StructToMap(response)

	// if err != nil {
	// 	return nil
	// }

	return response
}

func GetInv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	return os.Getenv(key)
}
