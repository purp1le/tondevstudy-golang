package main

import (
	"ton-lessons/internal/app"
	scan "ton-lessons/internal/scanner"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := app.InitApp(); err != nil {
		return err
	}

	scanner, err := scan.NewScanner()
	if err != nil {
		logrus.Error(err)
		return err
	}

	scanner.Listen()

	// client := liteclient.NewConnectionPool()
	// if err := client.AddConnectionsFromConfig(context.Background(), app.CFG.MAINNET_CONFIG); err != nil {
	// 	logrus.Error(err)
	// 	return err
	// }

	// uuid, err := uuid.NewUUID()
	// if err != nil {
	// 	logrus.Error(err)
	// 	return err
	// }

	// logrus.Info(uuid)
	// api := ton.NewAPIClient(client)

	// seed := app.CFG.Wallet.SEED
	// logrus.Info(seed)

	// wall, err := wallet.FromSeed(api, seed, wallet.HighloadV2Verified)
	// if err != nil {
	// 	logrus.Error(err)
	// 	return err
	// }

	// logrus.Info(wall.Address())

	// lastMaster, err := api.CurrentMasterchainInfo(context.Background())
	// if err != nil {
	// 	logrus.Error(err)
	// 	return err
	// }

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
	// 		var (
	// 			transferNotification structures.TransferNotification
	// 		)

	// 		if err := tlb.LoadFromCell(
	// 			&transferNotification,
	// 			inBody,
	// 			false,
	// 		); err != nil {
	// 			logrus.Error(err)
	// 			continue
	// 		}

	// 		logrus.Info("queryId: ", transferNotification.QueryId)
	// 		logrus.Info("jetton amount: ", transferNotification.Amount)
	// 		logrus.Info("sender address: ", transferNotification.Sender)
	// 	}
	// }

	return nil
}
