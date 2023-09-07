package application

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/config"
	"github.com/FaisalMashuri/my-wallet/infrastructure"
	accountRepository "github.com/FaisalMashuri/my-wallet/internal/domain/account/repository"
	notifRepository "github.com/FaisalMashuri/my-wallet/internal/domain/notification/repository"
	transactionRepository "github.com/FaisalMashuri/my-wallet/internal/domain/transaction/repository"

	notifController "github.com/FaisalMashuri/my-wallet/internal/domain/notification/controller"
	notifService "github.com/FaisalMashuri/my-wallet/internal/domain/notification/service"
	transactionController "github.com/FaisalMashuri/my-wallet/internal/domain/transaction/controller"
	transactionService "github.com/FaisalMashuri/my-wallet/internal/domain/transaction/service"

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

	//define repository
	userRepo := userRepository.NewRepository(db, log)
	accountRepo := accountRepository.NewRepository(db)
	transactionRepo := transactionRepository.NewRepository(db)
	notifRepo := notifRepository.NewRepository(db)
	//define service
	userSvc := userService.NewService(userRepo, log, accountRepo)
	transacetionSvc := transactionService.NewService(transactionRepo, accountRepo, notifRepo)
	notifSvc := notifService.NewService(notifRepo)

	//define controller
	userCtrl := userController.NewController(userSvc, log)
	transactionCtrl := transactionController.NewController(transacetionSvc)
	notifCtrl := notifController.NewController(notifSvc)

	//define route
	routeApp := router.NewRouter(router.RouteParams{
		UserController:        userCtrl,
		TransactionController: transactionCtrl,
		NotifController:       notifCtrl,
	})

	routeApp.SetupRoute(app)

	err = app.Listen(fmt.Sprintf(":%s", config.AppConfig.Port))
	if err != nil {
		log.Fatal("Failed to start application")
	}
}
