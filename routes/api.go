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

	router.Run(":4004")
}

// docker run -e MYSQL_HOST=127.0.0.1 -e MYSQL_USER=root -e MYSQL_PASSWORD=P@ssW0rd -e MYSQL_DBNAME=golang_startup -p 4004:4004 -d --name starupConDB my-startup-app

// * docker run --name mysqlDocker -e MYSQL_ROOT_PASSWORD=P@ssW0rd -e MYSQL_DATABASE=golang_startup -d -p 3333:3306 mysql:5.7
// docker run --name mysqlDocker -e MYSQL_ROOT_PASSWORD=P@ssW0rd -e MYSQL_DATABASE=golang_startup -d -p 3333:3306 mysql:5.7
