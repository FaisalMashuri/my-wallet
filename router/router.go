package router

import (
	"github.com/FaisalMashuri/my-wallet/config"
	midtrans_ext "github.com/FaisalMashuri/my-wallet/external/midtrans"
	"github.com/FaisalMashuri/my-wallet/internal/domain/account"
	"github.com/FaisalMashuri/my-wallet/internal/domain/mpin"
	"github.com/FaisalMashuri/my-wallet/internal/domain/notification"
	sseCtrl "github.com/FaisalMashuri/my-wallet/internal/domain/sse/controller"
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/middleware"
	"github.com/Saucon/errcntrct"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
	"log"
)

type RouteParams struct {
	UserController        user.UserController
	TransactionController transaction.TransactionController
	NotifController       notification.NotificationController
	NotifSseController    sseCtrl.NotificationSseController
	TopUpController       topup.TopUpController
	MidtransController    midtrans_ext.MidtransController
	AccountController     account.AccountController
	PinController         mpin.PinController
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
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:9999/swagger/oauth2-redirect.html",
	}))

	// Define routes with auth
	v1 := app.Group("/api/v1")
	v1.All("/midtrans/transfer-callback", r.RouteParams.MidtransController.PaymentHandlerNotification)

	v1.Route("/auth", func(router fiber.Router) {
		router.Post("/register", r.RouteParams.UserController.Register)
		router.Post("/login", r.RouteParams.UserController.Login)
		router.Use(middleware.NewAuthMiddleware(config.AppConfig.SecretKey))
		router.Get("/token", middleware.GetCredential)
	})
	v1.Post("/users/pin", r.RouteParams.PinController.CreatePin)

	v1.Use(middleware.NewAuthMiddleware(config.AppConfig.SecretKey))
	v1.Get("/users/detail", middleware.GetCredential, r.RouteParams.UserController.GetDetailUserJWT)
	v1.Post("/users/accounts", middleware.GetCredential, r.RouteParams.AccountController.CreateAccount)
	v1.Post("/topup/initialize", middleware.GetCredential, r.RouteParams.TopUpController.InitializeTopUp)
	v1.Post("/tranfer-inquiry", middleware.GetCredential, r.RouteParams.TransactionController.TransferInquiry)
	v1.Post("/transfer-exec", middleware.GetCredential, r.RouteParams.TransactionController.TransferExec)
	v1.Get("/notifications", middleware.GetCredential, r.RouteParams.NotifController.GetUserNotification)
	v1.Get("/notification-stream", middleware.GetCredential, r.RouteParams.NotifSseController.StreamNotifcation)
}
