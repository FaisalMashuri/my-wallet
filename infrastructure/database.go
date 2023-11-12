package infrastructure

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/config"
	domainAccount "github.com/FaisalMashuri/my-wallet/internal/domain/account"
	domainmPin "github.com/FaisalMashuri/my-wallet/internal/domain/mpin"
	domainNotification "github.com/FaisalMashuri/my-wallet/internal/domain/notification"
	domainTopup "github.com/FaisalMashuri/my-wallet/internal/domain/topup"
	domainTransaction "github.com/FaisalMashuri/my-wallet/internal/domain/transaction"
	domainUser "github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

type Database struct {
	DB    *gorm.DB
	Error error
}

var dbInstance *Database
var once sync.Once

func ConnectDB() *Database {
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			config.AppConfig.DatabaseConfig.Host,
			config.AppConfig.DatabaseConfig.Port,
			config.AppConfig.DatabaseConfig.User,
			config.AppConfig.DatabaseConfig.Password,
			config.AppConfig.DatabaseConfig.Name,
		)
		//define database logger
		dbLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		)

		dbConfig := gorm.Config{
			Logger: dbLogger,
		}

		db, err := gorm.Open(postgres.Open(dsn), &dbConfig)
		dbInstance = &Database{
			DB:    db,
			Error: err,
		}

		//
		//if os.Args[len(os.Args)-1] == "--migrate" {
		//	AutoMigrate(db)
		//} else if os.Args[len(os.Args)-1] == "--droptable" {
		//}
		AutoMigrate(dbInstance.DB)
	})

	return dbInstance
}

func AutoMigrate(db *gorm.DB) {
	fmt.Println("migrasi")
	db.Debug().AutoMigrate(
		domainUser.User{},
		domainAccount.Account{},
		domainTransaction.Transaction{},
		domainTransaction.TransactionInquiry{},
		domainNotification.Notification{},
		domainTopup.TopUp{},
		domainmPin.Pin{},
	)
}

func DropAllTable(db *gorm.DB) {
	fmt.Println("Drpoping Table")
	_ = db.Migrator().DropTable(
		domainUser.User{},
		domainAccount.Account{},
		domainTransaction.Transaction{},
		domainTransaction.TransactionInquiry{},
		domainNotification.Notification{},
		domainTopup.TopUp{},
		domainmPin.Pin{},
	)
}
