package scan

import (
	"context"
	"encoding/base64"
	"fmt"
	"ton-lessons/internal/app"
	"ton-lessons/internal/structures"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/tlb"
)

func (s *Scanner) processDedustSwap(
	master *tlb.BlockInfo,
	msgOut *tlb.ExternalMessageOut,
) error {
	var (
		swapEvent structures.DedustSwapEvent
	)

	if msgOut.Body == nil {
		return nil
	}

	if err := tlb.LoadFromCell(
		&swapEvent,
		msgOut.Body.BeginParse(),
		false,
	); err != nil {
		logrus.Error(err)
		return nil
	}

	accountInfo, err := s.Api.GetAccount(
		context.Background(),
		master,
		msgOut.SrcAddr,
	)
	if err != nil {
		return err
	}

	if base64.StdEncoding.EncodeToString(accountInfo.Code.Hash()) != app.DedustPoolCodeHash {
		logrus.Info("[DEDUST] code not valid")
		return nil
	}

	fmt.Println("[SCN] FOUND DEDUST SWAP")
	fmt.Printf("[SCN] ASSET IN - [%s] [%s]\n", swapEvent.AmountIn.TON(), swapEvent.AssetIn.Type())
	fmt.Printf("[SCN] ASSET OUT - [%s] [%s]\n", swapEvent.AmountOut.TON(), swapEvent.AssetOut.Type())
	fmt.Printf("[SCN] SENDER - [%s]\n", swapEvent.ExtraInfo.SenderAddr)
	fmt.Printf("[SCN] RESERVE 0 - [%s]\n", swapEvent.ExtraInfo.Reserve0.TON())
	fmt.Printf("[SCN] RESERVE 1 - [%s]\n", swapEvent.ExtraInfo.Reserve1.TON())

	return nil
}

func (s *Scanner) processDedustDeposit(
	master *tlb.BlockInfo,
	msgOut *tlb.ExternalMessageOut,
) error {
	var (
		deposit structures.DedustDepositEvent
	)

	if msgOut.Body == nil {
		return nil
	}

	accountInfo, err := s.Api.GetAccount(
		context.Background(),
		master,
		msgOut.SrcAddr,
	)
	if err != nil {
		return err
	}

	if base64.StdEncoding.EncodeToString(accountInfo.Code.Hash()) != app.DedustPoolCodeHash {
		logrus.Info("[DEDUST] code not valid")
		return nil
	}

	if err := tlb.LoadFromCell(
		&deposit,
		msgOut.Body.BeginParse(),
		false,
	); err != nil {
		return nil
	}

	fmt.Println("[SCN] found dedust deposit liq")
	fmt.Printf("[SCN] SENDER - [%s]\n", deposit.SenderAddr)
	fmt.Printf("[SCN] AMOUNT 0- [%s]\n", deposit.Amount0.String())
	fmt.Printf("[SCN] AMOUNT 1 - [%s]\n", deposit.Amount1.String())
	fmt.Printf("[SCN] RESERVE 0 - [%s]\n", deposit.Reserve0.String())
	fmt.Printf("[SCN] RESERVE 1 - [%s]\n", deposit.Reserve1.String())
	fmt.Printf("[SCN] LIQUIDITY - [%s]\n", deposit.Liquidity.String())

	return nil
}

func (s *Scanner) processDedustWithdrawal(
	master *tlb.BlockInfo,
	msgOut *tlb.ExternalMessageOut,
) error {
	var (
		withdraw structures.DedustDepositWithdrawal
	)

	if msgOut.Body == nil {
		return nil
	}

	accountInfo, err := s.Api.GetAccount(
		context.Background(),
		master,
		msgOut.SrcAddr,
	)
	if err != nil {
		return err
	}

	if base64.StdEncoding.EncodeToString(accountInfo.Code.Hash()) != app.DedustPoolCodeHash {
		logrus.Info("[DEDUST] code not valid")
		return nil
	}

	if err := tlb.LoadFromCell(
		&withdraw,
		msgOut.Body.BeginParse(),
		false,
	); err != nil {
		return nil
	}

	fmt.Println("[SCN] found dedust withdraw liq")
	fmt.Printf("[SCN] SENDER - [%s]\n", withdraw.SenderAddr)
	fmt.Printf("[SCN] AMOUNT 0- [%s]\n", withdraw.Amount0.String())
	fmt.Printf("[SCN] AMOUNT 1 - [%s]\n", withdraw.Amount1.String())
	fmt.Printf("[SCN] RESERVE 0 - [%s]\n", withdraw.Reserve0.String())
	fmt.Printf("[SCN] RESERVE 1 - [%s]\n", withdraw.Reserve1.String())
	fmt.Printf("[SCN] LIQUIDITY - [%s]\n", withdraw.Liquidity.String())

	return nil
}
