package routes

import (
	"bwastartup/config"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/users"
	"net/http"
	"runtime"
	"time"

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

	router.GET("/", func(ctx *gin.Context) {
		now := time.Now()
		nowFormat := now.Format("2006-01-02 15:04:05")

		map1 := map[string]string{
			"version":     runtime.Version(),
			"last_update": nowFormat,
		}

		// Convert the map to JSON
		jsonContent := helper.MapToJson(map1)
		config.Loggers("info", string(jsonContent))

		ctx.JSON(http.StatusOK, map1)
	})

	api := router.Group("/api/v1")
	api.POST("/users", userhandler.RegisterUser)
	router.Run(":3000")

	http.Serve(state.Listener, router)
}
