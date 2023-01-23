package helper

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
