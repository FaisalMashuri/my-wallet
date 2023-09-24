package config

import (
	"github.com/FaisalMashuri/my-wallet/shared/constant"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var AppConfig Config

type Config struct {
	Name           string `env:"NAME_SERVER"`
	Port           string `env:"PORT"`
	Host           string `env:"HOST_SERVER"`
	DatabaseConfig DatabaseConfig
	RabbitMQConfig RabbitMQ
	AppEnv         string `env:"APP_ENV"`
	SecretKey      string `env:"SECRET_KEY"`
	DSN_MQ         string `env:"DSN_MQ"`
	ErrorContract  ErrorContract
	Midtrans       MidtransConfig
}

type ErrorContract struct {
	JSONPathFile string `env:"PATH_ERR_CONTRACT"`
}

type DatabaseConfig struct {
	Name        string `env:"NAME_DB"`
	Host        string `env:"HOST_DB"`
	Port        string `env:"PORT_DB"`
	User        string `env:"USER_DB"`
	Password    string `env:"PASSWORD_DB"`
	DatabaseUrl string `env:"DATABASE_URL"`
}

type RabbitMQ struct {
	Host     string `env:"HOST_RABBIT"`
	Port     string `env:"PORT_RABBIT"`
	User     string `env:"USER_RABBIT"`
	Password string `env:"PASSWORD_RABBIT"`
}

type MidtransConfig struct {
	ServerKey  string `env:"MIDTRANS_SERVER_KEY"`
	ClientKey  string `env:"MIDTRANS_CLIENT_KEY"`
	MerchantID string `env:"MIDTRANS_MERCHANT_ID"`
	Env        string `env:"MIDTRANS_ENV"`
}

func LoadConfig() error {
	if os.Getenv("APP_ENV") != constant.AppEnvironmentProduction {
		err := godotenv.Load()
		if err != nil {
			log.Fatal()
		}
	}

	err := env.Parse(&AppConfig)
	if err != nil {
		log.Default()
		return err
	}

	err = env.Parse(&AppConfig.DatabaseConfig)
	if err != nil {
		log.Default()

		return err
	}

	err = env.Parse(&AppConfig.ErrorContract)
	if err != nil {
		log.Default()

		return err
	}

	err = env.Parse(&AppConfig.Midtrans)
	if err != nil {
		log.Default()
		return err
	}

	err = env.Parse(&AppConfig.RabbitMQConfig)
	if err != nil {
		return err
	}
	return err
}
