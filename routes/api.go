package routes

import (
	"startup/app/auth"
	"startup/app/campaign"
	"startup/app/handler"
	"startup/app/middleware"
	"startup/app/transaction"
	"startup/app/users"
	"startup/config"

	"github.com/gin-gonic/gin"
)

func InitApi() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies(nil)
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

	// transaction
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository, campaignRepository)
	transactionnhandler := handler.NewTransaction(transactionService)

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
	apiMiddleware.PUT("/campaigns/:id", campaignhandler.UpdateCampaign)
	apiMiddleware.POST("/campaigns-images", campaignhandler.UploadImage)

	apiMiddleware.GET("/campaigns/:id/transaction", transactionnhandler.GetCampaignTransaction)
	apiMiddleware.GET("/transaction", transactionnhandler.GetUserTransaction)

	router.Run(":4004")
}
