package handler

import (
	"bwastartup/config"
	"bwastartup/helper"
	"bwastartup/users"
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

	c.JSON(http.StatusOK, map1)
}

// <========== end version ==============>

type userHandler struct {
	userService users.Service
}

// <========== start register ==============>
// tangkap input dari user
// map input dari user ke struct RegisterUserInput
// struct di atas kita passing sebagai parameter service

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

	formatter := helper.FormatUser(user, "Tokentokentoken")
	response := helper.ApiResponse("account has been registered", http.StatusOK, "success", formatter)
	config.Loggers("info", response)

	c.JSON(http.StatusOK, response)
}

func NewUserHandler(userService users.Service) *userHandler {
	return &userHandler{userService}
}

// <========== end register ==============>

// <========== start login ==============>
// user memasukan input(email dan password)
// input ditangkap handler
// mapping dari input user ke input struct
// input struct passing service
// di service mencari dengan bantuan repository user dengan email dan mencocokan password

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

	formatter := helper.FormatUser(loggedinUser, "tokenasdad")
	response := helper.ApiResponse("successfully loggedin", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

// <========== end login ==============>
