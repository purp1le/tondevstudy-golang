package main

import (
	"context"
	"math/rand"
	"time"
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

	tonVaultAddr := address.MustParseAddr("EQDa4VOnTYlLvDJ0gZjNYm5PXfSmmtL6Vs6A_CZEtXCNICq_")
	// jettonVaultAddr := address.MustParseAddr("EQBeWd2_71HcPmAoTX2i9h0HWehA3_G76lxk90yyXmKXuje7")
	tonJettonPoolAddr := address.MustParseAddr("EQD0F_w35CTWUxTWRjefoV-400KRA2jX51X4ezIgmUUY_0Qn")

	dedustSwap := structures.DedustRequestNativeSwap{
		QueryId: rand.Uint64(),
		Amount:  tlb.MustFromTON("1"),
		SwapStep: structures.DedustSwapStep{
			PoolAddr: tonJettonPoolAddr,
			SwapStepParams: structures.DedustSwapStepParams{
				Limit: tlb.MustFromTON("0"),
				Next:  nil,
			},
		},
		SwapParams: structures.DedustSwapParams{
			Deadline:       uint32(time.Now().Unix()) + 60*60,
			RecipientAddr:  wall.Address(),
			ReferralAddr:   address.NewAddressNone(),
			FulfillPayload: nil,
			RejectPayload:  nil,
		},
	}

	swapBody, err := tlb.ToCell(&dedustSwap)
	if err != nil {
		return err
	}

	if err := wall.Send(
		context.Background(),
		wallet.SimpleMessage(
			tonVaultAddr,
			tlb.MustFromTON("1.3"),
			swapBody,
		), true,
	); err != nil {
		return err
	}

	// scanner, err := scan.NewScanner()
	// if err != nil {
	// 	logrus.Error(err)
	// 	return err
	// }

	// scanner.Listen()

	return nil
}
