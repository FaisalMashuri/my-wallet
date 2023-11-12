package application

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/config"
	midtransService "github.com/FaisalMashuri/my-wallet/external/midtrans/service"
	"github.com/FaisalMashuri/my-wallet/infrastructure"
	accountRepository "github.com/FaisalMashuri/my-wallet/internal/domain/account/repository"
	mPinController "github.com/FaisalMashuri/my-wallet/internal/domain/mpin/controller"
	mPinRepository "github.com/FaisalMashuri/my-wallet/internal/domain/mpin/repository"
	mPinService "github.com/FaisalMashuri/my-wallet/internal/domain/mpin/service"
	notifController "github.com/FaisalMashuri/my-wallet/internal/domain/notification/controller"
	"github.com/FaisalMashuri/my-wallet/internal/domain/notification/dto/response"
	notifRepository "github.com/FaisalMashuri/my-wallet/internal/domain/notification/repository"
	notifService "github.com/FaisalMashuri/my-wallet/internal/domain/notification/service"
	"github.com/FaisalMashuri/my-wallet/internal/domain/sse/controller"
	"github.com/FaisalMashuri/my-wallet/internal/domain/sse/dto"
	topupController "github.com/FaisalMashuri/my-wallet/internal/domain/topup/controller"
	topupRepository "github.com/FaisalMashuri/my-wallet/internal/domain/topup/repository"
	topupService "github.com/FaisalMashuri/my-wallet/internal/domain/topup/service"

	transactionController "github.com/FaisalMashuri/my-wallet/internal/domain/transaction/controller"
	transactionRepository "github.com/FaisalMashuri/my-wallet/internal/domain/transaction/repository"
	transactionService "github.com/FaisalMashuri/my-wallet/internal/domain/transaction/service"

	midtransController "github.com/FaisalMashuri/my-wallet/external/midtrans/controller"
	accountController "github.com/FaisalMashuri/my-wallet/internal/domain/account/controller"
	userController "github.com/FaisalMashuri/my-wallet/internal/domain/user/controller"

	accountService "github.com/FaisalMashuri/my-wallet/internal/domain/account/service"
	mqService "github.com/FaisalMashuri/my-wallet/internal/domain/mq/service"
	userRepository "github.com/FaisalMashuri/my-wallet/internal/domain/user/repository"
	userService "github.com/FaisalMashuri/my-wallet/internal/domain/user/service"

	"github.com/FaisalMashuri/my-wallet/middleware"
	"github.com/FaisalMashuri/my-wallet/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
)

func Run() {
	err := config.LoadConfig()
	log := infrastructure.Log

	if err != nil {
		os.Exit(1)
	}
	err = middleware.LoadErrorListFromJsonFile(config.AppConfig.ErrorContract.JSONPathFile)
	if err != nil {
		log.Error("Error connecting database")
	}

	database := infrastructure.ConnectDB()
	if database.Error != nil {
		log.Error("Error connecting database")
		os.Exit(1)
	}

	redisClient := infrastructure.RedisClient
	mq := infrastructure.NewRabbitMQ(config.AppConfig)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	app.Use(cors.New())
	hub := dto.Hub{NotificationChannel: map[string]chan response.NotificationDataRes{}}

	//define repository
	userRepo := userRepository.NewRepository(database.DB, log)
	accountRepo := accountRepository.NewRepository(database.DB)
	transactionRepo := transactionRepository.NewRepository(database.DB)
	notifRepo := notifRepository.NewRepository(database.DB)
	topupRepo := topupRepository.NewRepository(database.DB)
	mPinRepo := mPinRepository.NewRepository(database.DB)

	//define service
	messageQueueService := mqService.NewMqService(mq)
	userSvc := userService.NewService(userRepo, log, accountRepo, mPinRepo, redisClient, messageQueueService)
	transacetionSvc := transactionService.NewService(transactionRepo, accountRepo, notifRepo, &hub, redisClient)
	notifSvc := notifService.NewService(notifRepo)
	midtransSvc := midtransService.NewService()
	topUpSvc := topupService.NewService(topupRepo, midtransSvc, notifRepo, accountRepo, &hub)
	accountSvc := accountService.NewService(accountRepo)
	mPinSvc := mPinService.NewService(mPinRepo)

	//define controller
	userCtrl := userController.NewController(userSvc, log)
	transactionCtrl := transactionController.NewController(transacetionSvc, mPinSvc)
	notifCtrl := notifController.NewController(notifSvc)
	notifSseCtrl := controller.NewNotification(&hub)
	topUpCtrl := topupController.NewController(topUpSvc)
	midtransCtrl := midtransController.NewController(midtransSvc, topUpSvc)
	accountCtrl := accountController.NewController(accountSvc)
	mPinCtrl := mPinController.NewController(mPinSvc)
	//define route
	routeApp := router.NewRouter(router.RouteParams{
		UserController:        userCtrl,
		TransactionController: transactionCtrl,
		NotifController:       notifCtrl,
		NotifSseController:    notifSseCtrl,
		TopUpController:       topUpCtrl,
		MidtransController:    midtransCtrl,
		AccountController:     accountCtrl,
		PinController:         mPinCtrl,
	})

	routeApp.SetupRoute(app)

	err = app.Listen(fmt.Sprintf(":%s", config.AppConfig.Port))
	if err != nil {
		log.Fatal("Failed to start application")
	}
}
