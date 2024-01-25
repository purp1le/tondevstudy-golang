package app

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tvm/cell"
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
	START_BLOCK    uint64

	Wallet struct {
		SEED []string
	}
}

var CFG *appConfig = &appConfig{}
var DedustPoolCodeHash string
var dedustCode = "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"

func InitDedustCode() error {
	codeBytes, err := base64.StdEncoding.DecodeString(dedustCode)
	if err != nil {
		return err
	}

	codeCell, err := cell.FromBOC(codeBytes)
	if err != nil {
		logrus.Error(err)
		return err
	}

	DedustPoolCodeHash = base64.StdEncoding.EncodeToString(codeCell.Hash())
	return nil
}

func InitConfig() (err error) {

	if err := InitDedustCode(); err != nil {
		logrus.Error(err)
		return err
	}
	
	godotenv.Load(".env")

	CFG.Postgres.HOST = os.Getenv("POSTGRES_HOST")
	CFG.Postgres.PORT = os.Getenv("POSTGRES_PORT")
	CFG.Postgres.USER = os.Getenv("POSTGRES_USER")
	CFG.Postgres.PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	CFG.Postgres.NAME = os.Getenv("POSTGRES_DB")
	CFG.Postgres.SSLMODE = os.Getenv("POSTGRES_SSLMODE")
	CFG.Postgres.TIMEZONE = os.Getenv("POSTGRES_TIMEZONE")

	CFG.Logger.LOGLVL = os.Getenv("LOGL")

	CFG.START_BLOCK, err = strconv.ParseUint(os.Getenv("START_BLOCK"), 10, 64)
	if err != nil {
		return err
	}

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
