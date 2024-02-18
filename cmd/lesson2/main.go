package main

import (
	"ton-lessons/internal/app"
	scan "ton-lessons/internal/scanner"
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

	// liteclient := liteclient.NewConnectionPool()
	// if err := liteclient.AddConnectionsFromConfig(context.Background(), app.CFG.MAINNET_CONFIG); err != nil {
	// 	logrus.Error(err)
	// 	return err
	// }

	// api := ton.NewAPIClient(liteclient)

	// seed := app.CFG.Wallet.SEED
	// wall, err := wallet.FromSeed(api, seed, wallet.HighloadV2Verified)
	// if err != nil {
	// 	return err
	// }

	// uuid := uuid.New()
	// logrus.Info(uuid.String())

	// lastMaster, err := api.CurrentMasterchainInfo(context.Background())
	// if err != nil {
	// 	return err
	// }

	// accountInfo, err := api.GetAccount(
	// 	context.Background(),
	// 	lastMaster,
	// 	wall.Address(),
	// )
	// if err != nil {
	// 	return err
	// }

	// transactions := make(chan *tlb.Transaction)

	// go api.SubscribeOnTransactions(
	// 	context.Background(),
	// 	wall.Address(),
	// 	accountInfo.LastTxLT,
	// 	transactions,
	// )

	// jettonWalletAddress := "EQCaAgX3aSw3P7ZBklcRvMFlhBcCnF9HPuTOtfL9fZaZe1YL"

	// for {
	// 	select {
	// 	case newTransaction := <-transactions:
	// 		if newTransaction.IO.In.MsgType != tlb.MsgTypeInternal {
	// 			logrus.Info("not internal message")
	// 			continue
	// 		}

	// 		internalMessage := newTransaction.IO.In.AsInternal()
	// 		if internalMessage.SrcAddr.String() != jettonWalletAddress {
	// 			continue
	// 		}

	// 		if internalMessage.Body == nil {
	// 			logrus.Info("internal message body = nil")
	// 			continue
	// 		}

	// 		slice := internalMessage.Body.BeginParse()

	// 		var transferNotification structures.JettonNotification

	// 		if err := tlb.LoadFromCell(
	// 			&transferNotification,
	// 			slice,
	// 			false,
	// 		); err != nil {
	// 			logrus.Error(err)
	// 			continue
	// 		}

	// 		logrus.Info("found jetton transfer!")
	// 		logrus.Info("Sender: ", transferNotification.Sender)
	// 		logrus.Info("Amount: ", transferNotification.Amount.String())
	// 		logrus.Info("QueryId: ", transferNotification.QueryId)

	// 	}
	// }

	scanner, err := scan.NewScanner()
	if err != nil {
		return err
	}

	scanner.Listen()

	return nil
}
