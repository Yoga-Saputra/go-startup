package routes

import (
	"startup/app/auth"
	"startup/app/campaign"
	"startup/app/handler"
	"startup/app/middleware"
	"startup/app/users"
	"startup/config"

	"github.com/gin-gonic/gin"
)

func InitApi() {
	router := gin.Default()
	router.Use(gin.Logger())

	db := config.ConnectionDB()

	// user
	userRepository := users.NewRepository(db)
	userService := users.NewService(userRepository)
	authService := auth.NewService()
	userhandler := handler.NewUserHandler(userService, authService)

	// campaign
	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignhandler := handler.NewCampaignHandler(campaignService)

	router.Static("/images", "./storage/images")
	router.GET("/", handler.Version)

	// routes prefix
	api := router.Group("/api/v1")
	api.POST("/users", userhandler.RegisterUser)
	api.POST("/session", userhandler.Login)
	api.POST("/email_checkers", userhandler.CheckEmailAvailability)
	api.GET("/campaigns", campaignhandler.GetAllCamp)
	api.GET("/campaigns/:id", campaignhandler.GetCampain)

	// middleware grouping
	apiMiddleware := api.Use(middleware.AuthMiddleware(authService, userService))
	apiMiddleware.POST("/avatars", userhandler.UploadAvatar)
	apiMiddleware.POST("/campaigns", campaignhandler.CreateCampaign)

	router.Run(":3000")
}
