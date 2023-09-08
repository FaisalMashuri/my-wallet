package router

import (
	"github.com/FaisalMashuri/my-wallet/config"
	"github.com/FaisalMashuri/my-wallet/internal/domain/notification"
	sseCtrl "github.com/FaisalMashuri/my-wallet/internal/domain/sse/controller"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/middleware"
	"github.com/Saucon/errcntrct"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"log"
)

type RouteParams struct {
	UserController        user.UserController
	TransactionController transaction.TransactionController
	NotifController       notification.NotificationController
	NotifSseController    sseCtrl.NotificationSseController
}

type router struct {
	RouteParams RouteParams
	Log         logrus.Logger
}

func NewRouter(params RouteParams) router {
	return router{RouteParams: params}
}

func (r *router) SetupRoute(app *fiber.App) {
	if err := errcntrct.InitContract(config.AppConfig.ErrorContract.JSONPathFile); err != nil {
		//logger.Fatal(err, "main : init contract", nil)
		log.Fatal("main : init contract " + err.Error())
	}

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON("HALO")
	})

	// Define routes with auth
	v1 := app.Group("/api/v1")

	v1.Route("/auth", func(router fiber.Router) {
		router.Post("/register", r.RouteParams.UserController.Register)
		router.Post("/login", r.RouteParams.UserController.Login)
		router.Use(middleware.NewAuthMiddleware(config.AppConfig.SecretKey))
		router.Get("/token", middleware.GetCredential)
	})

	v1.Use(middleware.NewAuthMiddleware(config.AppConfig.SecretKey))
	v1.Post("/tranfer-inquiry", middleware.GetCredential, r.RouteParams.TransactionController.TransferInquiry)
	v1.Post("/transfer-exec", middleware.GetCredential, r.RouteParams.TransactionController.TransferExec)
	v1.Get("/notifications", middleware.GetCredential, r.RouteParams.NotifController.GetUserNotification)
	v1.Get("/notification-stream", middleware.GetCredential, r.RouteParams.NotifSseController.StreamNotifcation)
}
