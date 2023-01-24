package handler

import (
	"bwastartup/app/auth"
	"bwastartup/app/helper"
	"bwastartup/app/users"
	"bwastartup/config"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
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
	var input users.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		helper.ErrorValidation(err, c, "register account failed", "error", http.StatusUnprocessableEntity, errors)
		return
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
	if err != nil {
		errorMsg := gin.H{"error": "server error"}
		helper.ErrorValidation(err, c, "email checking failed", "error", http.StatusBadRequest, errorMsg)
		return
	}

	data := gin.H{"is_available": isEmailAvailable}

	var metaMsg string
	metaMsg = "email has been registerd"
	if isEmailAvailable {
		metaMsg = "email is available"
	}
	response := helper.ApiResponse(metaMsg, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

// <========== end check email availability ==============>

// <========== start upload avatar ==============>
func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	responseErrorUploadAvatar(err, c)
	userId := 10
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
	return token
}
