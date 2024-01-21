package app

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/liteclient"
)

type appConfig struct {
	Postgres struct {
		HOST, PORT, USER, PASSWORD,
		NAME, SSLMODE, TIMEZONE string
	}

	Logger struct {
		LOGLVL string
	}

	MAINNET_CONFIG *liteclient.GlobalConfig

	Wallet struct {
		SEED []string
	}
}

var CFG *appConfig = &appConfig{}

func InitConfig() (err error) {
	godotenv.Load(".env")

	CFG.Postgres.HOST = os.Getenv("POSTGRES_HOST")
	CFG.Postgres.PORT = os.Getenv("POSTGRES_PORT")
	CFG.Postgres.USER = os.Getenv("POSTGRES_USER")
	CFG.Postgres.PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	CFG.Postgres.NAME = os.Getenv("POSTGRES_DB")
	CFG.Postgres.SSLMODE = os.Getenv("POSTGRES_SSLMODE")
	CFG.Postgres.TIMEZONE = os.Getenv("POSTGRES_TIMEZONE")

	CFG.Logger.LOGLVL = os.Getenv("LOGL")

	jsonConfig, err := os.Open("mainnet-config.json")
	if err == nil {

		if err := json.NewDecoder(jsonConfig).Decode(&CFG.MAINNET_CONFIG); err != nil {
			return err
		}
	} else {
		logrus.Error(err)
		CFG.MAINNET_CONFIG = nil
	}
	defer jsonConfig.Close()

	CFG.Wallet.SEED = strings.Split(os.Getenv("SEED"), ";")

	return nil
}
