package main

import (
	"context"
	"math/big"
	"ton-lessons/internal/app"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
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

	client := liteclient.NewConnectionPool()
	if err := client.AddConnectionsFromConfig(context.Background(), app.CFG.MAINNET_CONFIG); err != nil {
		logrus.Error(err)
		return err
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Info(uuid.String())

	api := ton.NewAPIClient(client)

	seed := app.CFG.Wallet.SEED
	logrus.Info(seed)

	wall, err := wallet.FromSeed(api, seed, wallet.HighloadV2Verified)
	if err != nil {
		logrus.Error(err)
		return err
	}

	lastMaster, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		logrus.Error(err)
		return err
	}

	jettonWalletAddress := address.MustParseAddr("EQCaAgX3aSw3P7ZBklcRvMFlhBcCnF9HPuTOtfL9fZaZe1YL")

	execResult, err := api.RunGetMethod(context.Background(), lastMaster, jettonWalletAddress, "get_wallet_data")
	if err != nil {
		logrus.Error(err)
		return err
	}

	balance := execResult.AsTuple()[0].(*big.Int)
	logrus.Info(balance)

	if err := wall.Send(
		context.Background(),
		wallet.SimpleMessage(
			jettonWalletAddress,
			tlb.MustFromTON("0.3"),
			cell.BeginCell().
				MustStoreUInt(0x0f8a7ea5, 32).
				MustStoreUInt(0, 64).
				MustStoreBigCoins(balance).
				MustStoreAddr(address.MustParseAddr("EQDRCfK1eHMRQvZEg0Ylb3AzGymGnsu1TQnELGiOQZx7M6KO")).
				MustStoreAddr(address.MustParseAddr("EQDRCfK1eHMRQvZEg0Ylb3AzGymGnsu1TQnELGiOQZx7M6KO")).
				MustStoreUInt(0, 1).
				MustStoreCoins(100000000).
				MustStoreUInt(0, 1).
				EndCell(),
		),
		true,
	); err != nil {
		logrus.Error(err)
		return err
	}
	// accountInfo, err := api.GetAccount(context.Background(), lastMaster, wall.Address())
	// if err != nil {
	// 	logrus.Error(err)
	// 	return err
	// }

	// transactions := make(chan *tlb.Transaction)
	// go api.SubscribeOnTransactions(context.Background(), wall.Address(), accountInfo.LastTxLT, transactions)
	// for {
	// 	select {
	// 	case newTransaction := <-transactions:
	// 		if newTransaction.IO.In.MsgType != tlb.MsgTypeInternal {
	// 			continue
	// 		}

	// 		internalMessage := newTransaction.IO.In.AsInternal()
	// 		logrus.Info("sender address: ", internalMessage.SrcAddr)
	// 		logrus.Info("receiver address: ", internalMessage.DstAddr)
	// 		logrus.Info("amount: ", internalMessage.Amount)

	// 		if internalMessage.Body == nil {
	// 			continue
	// 		}

	// 		inBody := internalMessage.Body.BeginParse()
	// 		opcode, err := inBody.LoadUInt(32)
	// 		if err != nil {
	// 			logrus.Error(err)
	// 			continue
	// 		}

	// 		logrus.Info(opcode)
	// 		if opcode != 0 {
	// 			continue
	// 		}

	// 		comment, err := inBody.LoadStringSnake()
	// 		if err != nil {
	// 			logrus.Error(err)
	// 			continue
	// 		}

	// 		logrus.Info("comment: ", comment)
	// 		if comment == uuid.String() {
	// 			logrus.Info("deposit confirmed, uuid equals")
	// 		}
	// 	}
	// }

	return nil
}
