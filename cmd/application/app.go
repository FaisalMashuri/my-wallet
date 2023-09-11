package application

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/config"
	midtransService "github.com/FaisalMashuri/my-wallet/external/midtrans/service"
	"github.com/FaisalMashuri/my-wallet/infrastructure"
	accountRepository "github.com/FaisalMashuri/my-wallet/internal/domain/account/repository"
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
	userController "github.com/FaisalMashuri/my-wallet/internal/domain/user/controller"

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
	if err != nil {
		os.Exit(1)
	}

	log := infrastructure.Log
	db, err := infrastructure.ConnectDB()
	if err != nil {
		log.Error("Error connecting database")
	}
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return middleware.NewErrorhandler(ctx, err)
		},
	})
	app.Use(cors.New())
	hub := dto.Hub{NotificationChannel: map[string]chan response.NotificationDataRes{}}

	//define repository
	userRepo := userRepository.NewRepository(db, log)
	accountRepo := accountRepository.NewRepository(db)
	transactionRepo := transactionRepository.NewRepository(db)
	notifRepo := notifRepository.NewRepository(db)
	topupRepo := topupRepository.NewRepository(db)

	//define service
	userSvc := userService.NewService(userRepo, log, accountRepo)
	transacetionSvc := transactionService.NewService(transactionRepo, accountRepo, notifRepo, &hub)
	notifSvc := notifService.NewService(notifRepo)
	midtransSvc := midtransService.NewService()
	topUpSvc := topupService.NewService(topupRepo, midtransSvc, notifRepo, accountRepo, &hub)

	//define controller
	userCtrl := userController.NewController(userSvc, log)
	transactionCtrl := transactionController.NewController(transacetionSvc)
	notifCtrl := notifController.NewController(notifSvc)
	notifSseCtrl := controller.NewNotification(&hub)
	topUpCtrl := topupController.NewController(topUpSvc)
	midtransCtrl := midtransController.NewController(midtransSvc, topUpSvc)
	//define route
	routeApp := router.NewRouter(router.RouteParams{
		UserController:        userCtrl,
		TransactionController: transactionCtrl,
		NotifController:       notifCtrl,
		NotifSseController:    notifSseCtrl,
		TopUpController:       topUpCtrl,
		MidtransController:    midtransCtrl,
	})

	routeApp.SetupRoute(app)

	fmt.Println("MIDTRANS SERVER KEY : ", config.AppConfig.Midtrans.ServerKey)

	fmt.Println("MIDTRANS Client KEY : ", config.AppConfig.Midtrans.ClientKey)
	fmt.Println("MIDTRANS Merchant ID : ", config.AppConfig.Midtrans.MerchantID)

	err = app.Listen(fmt.Sprintf(":%s", config.AppConfig.Port))
	if err != nil {
		log.Fatal("Failed to start application")
	}
}
