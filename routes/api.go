package routes

import (
	"net/http"
	"startup/app/auth"
	"startup/app/campaign"
	"startup/app/handler"
	"startup/app/middleware"
	"startup/app/users"
	"startup/config"

	"github.com/gin-gonic/gin"
	"github.com/jpillora/overseer"
)

func InitApi(state overseer.State) {
	router := gin.Default()
	router.Use(gin.Logger())

	db := config.ConnectionDB()
	userRepository := users.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	userService := users.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()

	userhandler := handler.NewUserHandler(userService, authService)
	campaignhandler := handler.NewCampaignHandler(campaignService)

	router.Static("/images", "./storage/images")
	router.GET("/", handler.Version)

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

	http.Serve(state.Listener, router)
}
