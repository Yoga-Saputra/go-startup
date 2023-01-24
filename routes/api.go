package routes

import (
	"bwastartup/app/auth"
	"bwastartup/app/handler"
	"bwastartup/app/users"
	"bwastartup/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jpillora/overseer"
)

func InitApi(state overseer.State) {
	router := gin.Default()
	router.Use(gin.Logger())

	db := config.ConnectionDB()
	userRepository := users.NewRepository(db)
	userService := users.NewService(userRepository)
	apiService := auth.NewService()
	userhandler := handler.NewUserHandler(userService, apiService)

	router.GET("/", handler.Version)

	api := router.Group("/api/v1")
	api.POST("/users", userhandler.RegisterUser)
	api.POST("/session", userhandler.Login)
	api.POST("/email_checkers", userhandler.CheckEmailAvailability)
	api.POST("/avatars", userhandler.UploadAvatar)

	router.Run(":3000")

	http.Serve(state.Listener, router)
}
