package handler

import (
	"fmt"
	"net/http"
	"runtime"
	"startup/app/auth"
	"startup/app/helper"
	"startup/app/users"
	"startup/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// <========== start version ==============>

func Version(c *gin.Context) {
	now := time.Now()
	nowFormat := now.Format("2006-01-02 15:04:05")
	map1 := map[string]string{
		"last_update": nowFormat,
		"build_with":  "GO",
		"version":     runtime.Version(),
	}

	// Convert the map to JSON
	jsonContent := helper.MapToJson(map1)
	config.Loggers("info", string(jsonContent))

	fmt.Println()

	c.JSON(http.StatusOK, map1)
}

// <========== end version ==============>

type userHandler struct {
	userService users.Service
	authService auth.Service
}

func NewUserHandler(userService users.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

// <========== start register ==============>
func (h *userHandler) RegisterUser(c *gin.Context) {
	var inputCheckEmail users.CheckEmailInput
	err := c.ShouldBindBodyWith(&inputCheckEmail, binding.JSON)
	if err != nil {
		errors := helper.FormatValidationError(err)
		helper.ErrorValidation(err, c, "email checking failed", "error", http.StatusUnprocessableEntity, errors)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(inputCheckEmail)
	if err != nil {
		errorMsg := gin.H{"error": "server error"}
		helper.ErrorValidation(err, c, "email checking failed", "error", http.StatusBadRequest, errorMsg)
		return
	}

	if !isEmailAvailable {
		data := gin.H{"is_available": isEmailAvailable}
		metaMsg := "email has been registerd"
		response := helper.ApiResponse(metaMsg, http.StatusOK, "error", data)
		c.JSON(http.StatusForbidden, response)
		return
	}

	var input users.RegisterUserInput
	err = c.ShouldBindBodyWith(&input, binding.JSON)
	if err != nil {
		errors := helper.FormatValidationError(err)
		helper.ErrorValidation(err, c, "register account failed", "error", http.StatusUnprocessableEntity, errors)
		// return
		c.JSON(http.StatusOK, err)
	}

	user, err := h.userService.RegisterUser(input)

	if err != nil {
		helper.ErrorValidation(err, c, "register account failed", "error", http.StatusBadRequest, err.Error())
		return
	}

	token := responseToken(user.ID, h, c, "register account failed")

	formatter := helper.FormatUser(user, token)
	response := helper.ApiResponse("account has been registered", http.StatusOK, "success", formatter)
	config.Loggers("info", response)

	c.JSON(http.StatusOK, response)
}

// <========== end register ==============>

// <========== start login ==============>
func (h *userHandler) Login(c *gin.Context) {
	var input users.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		helper.ErrorValidation(err, c, "login account failed", "error", http.StatusUnprocessableEntity, errors)
		return
	}

	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		helper.ErrorValidation(err, c, "login account failed", "error", http.StatusBadRequest, err.Error())
		return
	}

	token := responseToken(loggedinUser.ID, h, c, "login failed")
	formatter := helper.FormatUser(loggedinUser, token)
	response := helper.ApiResponse("successfully loggedin", http.StatusOK, "success", formatter)
	config.Loggers("info", response)
	c.JSON(http.StatusOK, response)
}

// <========== end login ==============>

// <========== start check email availability ==============>
func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input users.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		helper.ErrorValidation(err, c, "email checking failed", "error", http.StatusUnprocessableEntity, errors)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	config.Loggers("error", isEmailAvailable)
	if err != nil {
		errorMsg := gin.H{"error": "server error"}
		helper.ErrorValidation(err, c, "email checking failed", "error", http.StatusBadRequest, errorMsg)
		return
	}

	data := gin.H{"is_available": isEmailAvailable}

	metaMsg := "email has been registerd"
	success := "error"
	code := http.StatusForbidden

	if isEmailAvailable {
		code = http.StatusOK
		success = "success"
		metaMsg = "email is available"
	}
	response := helper.ApiResponse(metaMsg, code, success, data)
	c.JSON(code, response)
}

// <========== end check email availability ==============>

// <========== start upload avatar ==============>
func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	responseErrorUploadAvatar(err, c)
	currentUser := c.MustGet("current_user").(users.User)
	userId := currentUser.ID
	// sprintf = menggabungkan string
	path := fmt.Sprintf("storage/images/users/%d-%s", userId, file.Filename)
	err = c.SaveUploadedFile(file, path)
	responseErrorUploadAvatar(err, c)

	_, err = h.userService.SaveAvatar(userId, path)
	responseErrorUploadAvatar(err, c)
	responseSuccessUploadAvatar(c)
}

// <========== end upload avatar ==============>

func responseErrorUploadAvatar(err error, c *gin.Context) {
	if err != nil {
		errorMsg := gin.H{"is_uploaded": false}
		jsonContent := helper.InterfaceToJson(errorMsg)
		response := helper.ApiResponse("failed to upload avatar", http.StatusBadRequest, "error", errorMsg)
		c.JSON(http.StatusBadRequest, response)
		config.Loggers("error", string(jsonContent))
		return
	}
}

func responseSuccessUploadAvatar(c *gin.Context) {
	successMsg := gin.H{"is_uploaded": true}
	response := helper.ApiResponse("avatar successfully uploaded", http.StatusOK, "success", successMsg)
	config.Loggers("info", response)
	c.JSON(http.StatusOK, response)
}

func responseToken(id int, h *userHandler, c *gin.Context, msg string) string {
	token, err := h.authService.GenerateToken(id)
	if err != nil {
		helper.ErrorValidation(err, c, msg, "error", http.StatusBadRequest, err)
		return err.Error()
	}
	mapToken := map[string]interface{}{
		"user_id": id,
		"token":   token,
	}
	// Convert the map interface to JSON
	jsonToken := helper.InterfaceToJson(mapToken)
	config.Loggers("info", string(jsonToken))
	return token
}
