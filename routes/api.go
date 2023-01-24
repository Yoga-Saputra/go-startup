package routes

import (
	"bwastartup/config"
	"bwastartup/handler"
	"bwastartup/users"
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
	userhandler := handler.NewUserHandler(userService)

	router.GET("/", handler.Version)

	// start register and login routes

	api := router.Group("/api/v1")
	api.POST("/users", userhandler.RegisterUser)
	api.POST("/session", userhandler.Login)
	api.POST("/email_checkers", userhandler.CheckEmailAvailability)
	api.POST("/avatars", userhandler.UploadAvatar)

	// end register and login routes

	router.Run(":3000")

	http.Serve(state.Listener, router)
}
