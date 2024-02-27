package main

import (
	"context"
	"math/rand"
	"ton-lessons/internal/app"
	"ton-lessons/internal/structures"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/address"
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

	stonfiJettonRouterWallet := address.MustParseAddr("EQARULUYsmJq1RiZ-YiH-IJLcAZUVkVff-KBPwEmmaQGH6aC")
	// ptonJettonWalletAddr := address.MustParseAddr("EQARULUYsmJq1RiZ-YiH-IJLcAZUVkVff-KBPwEmmaQGH6aC")
	JetTonJettonWalletAddr := address.MustParseAddr("EQCaAgX3aSw3P7ZBklcRvMFlhBcCnF9HPuTOtfL9fZaZe1YL")
	stonfiRouterAddr := address.MustParseAddr("EQB3ncyBUTjZUA5EnFKR5_EnOMI9V1tTEAAPaiU71gc4TiUt")

	stonfiSwapBody := structures.StonfiSwapRequest{
		TokenWallet1: stonfiJettonRouterWallet,
		MinOut:       tlb.MustFromTON("0"),
		ToAddress:    wall.Address(),
		HasRef:       false,
		RefAddress:   nil,
	}

	stonfiSwapBodyCell, err := tlb.ToCell(&stonfiSwapBody)
	if err != nil {
		return err
	}

	transferRequest := structures.JettonTrasfer{
		QueryId:             rand.Uint64(),
		Amount:              tlb.MustFromTON("1.8"),
		Destination:         stonfiRouterAddr,
		ResponseDestination: stonfiRouterAddr,
		CustomPayload:       nil,
		FwdTonAmount:        tlb.MustFromTON("0.3"),
		FwdPayload:          stonfiSwapBodyCell,
	}

	transferRequestCell, err := tlb.ToCell(transferRequest)
	if err != nil {
		return err
	}

	if err := wall.Send(
		context.Background(),
		wallet.SimpleMessage(
			JetTonJettonWalletAddr,
			tlb.MustFromTON("0.5"),
			transferRequestCell,
		),
		true,
	); err != nil {
		return err
	}

	// ton -> jetton -> scale

	// tonVaultAddr := address.MustParseAddr("EQDa4VOnTYlLvDJ0gZjNYm5PXfSmmtL6Vs6A_CZEtXCNICq_")
	// // jettonVaultAddr := address.MustParseAddr("EQBeWd2_71HcPmAoTX2i9h0HWehA3_G76lxk90yyXmKXuje7")
	// // jettonWalletAddr := address.MustParseAddr("EQCaAgX3aSw3P7ZBklcRvMFlhBcCnF9HPuTOtfL9fZaZe1YL")
	// tonJettonPoolAddr := address.MustParseAddr("EQD0F_w35CTWUxTWRjefoV-400KRA2jX51X4ezIgmUUY_0Qn")
	// jettonScalePoolAddr := address.MustParseAddr("EQCVYWRk1gM3pjP8T7zVZYi7SViD0y_zOQ7YUbnU6u1f44tQ")

	// dedustSwap := structures.DedustRequestNativeSwap{
	// 	QueryId: rand.Uint64(),
	// 	Amount:  tlb.MustFromTON("1.5"),
	// 	SwapStep: structures.DedustSwapStep{
	// 		PoolAddr: tonJettonPoolAddr,
	// 		SwapStepParams: structures.DedustSwapStepParams{
	// 			Limit: tlb.MustFromTON("0"),
	// 			Next: &structures.DedustSwapStep{
	// 				PoolAddr: jettonScalePoolAddr,
	// 				SwapStepParams: structures.DedustSwapStepParams{
	// 					Limit: tlb.MustFromTON("0"),
	// 					Next:  nil,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	SwapParams: structures.DedustSwapParams{
	// 		Deadline:       uint32(time.Now().Unix()) + 60*60,
	// 		RecipientAddr:  wall.Address(),
	// 		ReferralAddr:   address.NewAddressNone(),
	// 		FulfillPayload: nil,
	// 		RejectPayload:  nil,
	// 	},
	// }

	// swapBody, err := tlb.ToCell(&dedustSwap)
	// if err != nil {
	// 	return err
	// }

	// if err := wall.Send(
	// 	context.Background(),
	// 	wallet.SimpleMessage(
	// 		tonVaultAddr,
	// 		tlb.MustFromTON("2"),
	// 		swapBody,
	// 	), true,
	// ); err != nil {
	// 	return err
	// }

	// scanner, err := scan.NewScanner()
	// if err != nil {
	// 	logrus.Error(err)
	// 	return err
	// }

	// scanner.Listen()

	return nil
}
