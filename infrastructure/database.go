package infrastructure

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/config"
	domainAccount "github.com/FaisalMashuri/my-wallet/internal/domain/account"
	domainNotification "github.com/FaisalMashuri/my-wallet/internal/domain/notification"
	domainTransaction "github.com/FaisalMashuri/my-wallet/internal/domain/transaction"
	domainUser "github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func ConnectDB() (*gorm.DB, error) {
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
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	dbConfig := gorm.Config{
		Logger: dbLogger,
	}

	db, err := gorm.Open(postgres.Open(dsn), &dbConfig)
	if err != nil {
		return nil, err
	}

	fmt.Println(os.Args[len(os.Args)-1])
	if os.Args[len(os.Args)-1] == "--migrate" {
		AutoMigrate(db)
	} else if os.Args[len(os.Args)-1] == "--droptable" {
		DropAllTable(db)
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB) {
	fmt.Println("migrasi")
	db.Debug().AutoMigrate(
		domainUser.User{},
		domainAccount.Account{},
		domainTransaction.Transaction{},
		domainTransaction.TransactionInquiry{},
		domainNotification.Notification{},
	)
}

func DropAllTable(db *gorm.DB) {
	fmt.Println("Drpoping Table")
	db.Migrator().DropTable(
		domainUser.User{},
		domainAccount.Account{},
		domainTransaction.Transaction{},
		domainTransaction.TransactionInquiry{},
		domainNotification.Notification{},
	)
}
