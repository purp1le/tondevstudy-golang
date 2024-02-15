package main

import (
	"context"
	"ton-lessons/internal/app"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := app.InitApp(); err != nil {
		logrus.Error(err)
		return err
	}

	liteclient := liteclient.NewConnectionPool()
	if err := liteclient.AddConnectionsFromConfig(context.Background(), app.CFG.MAINNET_CONFIG); err != nil {
		logrus.Error(err)
		return err
	}

	api := ton.NewAPIClient(liteclient)

	seed := app.CFG.Wallet.SEED
	wall, err := wallet.FromSeed(api, seed, wallet.HighloadV2Verified)
	if err != nil {
		return err
	}

	uuid := uuid.New()
	logrus.Info(uuid.String())

	lastMaster, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return err
	}

	accountInfo, err := api.GetAccount(
		context.Background(),
		lastMaster,
		wall.Address(),
	)
	if err != nil {
		return err
	}

	transactions := make(chan *tlb.Transaction)

	go api.SubscribeOnTransactions(
		context.Background(),
		wall.Address(),
		accountInfo.LastTxLT,
		transactions,
	)

	for {
		select {
		case newTransaction := <-transactions:
			if newTransaction.IO.In.MsgType != tlb.MsgTypeInternal {
				logrus.Info("not internal message")
				continue
			}

			internalMessage := newTransaction.IO.In.AsInternal()
			if internalMessage.Body == nil {
				logrus.Info("internal message body = nil")
				continue
			}

			slice := internalMessage.Body.BeginParse()
			opcode, err := slice.LoadUInt(32)
			if err != nil {
				logrus.Info("no have opcode")
				continue
			}

			if opcode != 0 {
				logrus.Info("not text message")
				continue
			}

			msg, err := slice.LoadStringSnake()
			if err != nil {
				logrus.Info("load string snake error")
				continue
			}

			logrus.Info("MSG: ", msg)
			logrus.Info("TON AMOUNT: ", internalMessage.Amount)
			logrus.Info("SENDER ADDRESS: ", internalMessage.SrcAddr)

			if msg == uuid.String() {
				logrus.Info("deposit successful")
			} else {
				logrus.Info("user not found")
			}
		}
	}

	return nil
}
