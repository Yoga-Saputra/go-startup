package helper

import "crypto/md5"

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

func ApiResponse(message string, code int, status string, data interface{}) response {
	meta := meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	response := response{
		Meta: meta,
		Data: data,
	}

	return response
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
