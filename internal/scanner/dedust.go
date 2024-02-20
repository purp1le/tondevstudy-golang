package scan

import (
	"context"
	"encoding/base64"
	"fmt"
	"ton-lessons/internal/app"
	"ton-lessons/internal/structures"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
)

func (s *Scanner) processDedustSwap(
	master *tlb.BlockInfo,
	msgOut *tlb.ExternalMessageOut,
) error {
	var (
		err        error
		opcode     uint64
		assetIn    string
		assetOut   string
		amountIn   float64
		amountOut  float64
		senderAddr string
		// refAddr    string
		reserve0 float64
		reserve1 float64
	)

	if msgOut.Body == nil {
		return nil
	}

	bodySl := msgOut.Body.BeginParse()

	opcode, err = bodySl.LoadUInt(32)
	if err != nil {
		return nil
	}

	if opcode != 0x9c610de3 {
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

	assetInType, err := bodySl.LoadUInt(4)
	if err != nil {
		return nil
	}

	if assetInType == 0 {
		assetIn = "TON"
	} else {
		workchain, err := bodySl.LoadUInt(8)
		if err != nil {
			return nil
		}

		addrBytes, err := bodySl.LoadSlice(256)
		if err != nil {
			return nil
		}

		assetIn = address.NewAddress(0, byte(workchain), addrBytes).String()
	}

	assetOutType, err := bodySl.LoadUInt(4)
	if err != nil {
		return nil
	}

	if assetOutType == 0 {
		assetOut = "TON"
	} else {
		workchain, err := bodySl.LoadUInt(8)
		if err != nil {
			return nil
		}

		addrBytes, err := bodySl.LoadSlice(256)
		if err != nil {
			return nil
		}

		assetOut = address.NewAddress(0, byte(workchain), addrBytes).String()
	}

	amountOutBig, err := bodySl.LoadBigCoins()
	if err != nil {
		return nil
	}

	amountOut, _ = amountOutBig.Float64()
	amountOut /= 1e9

	amountInBig, err := bodySl.LoadBigCoins()
	if err != nil {
		return nil
	}

	amountIn, _ = amountInBig.Float64()
	amountIn /= 1e9

	nextSl, err := bodySl.LoadRef()
	if err != nil {
		return nil
	}

	senderAddrType, err := nextSl.LoadAddr()
	if err != nil {
		return nil
	}

	senderAddr = senderAddrType.String()

	_, err = nextSl.LoadAddr()
	if err != nil {
		return nil
	}

	reserve0Big, err := nextSl.LoadBigCoins()
	if err != nil {
		return nil
	}

	reserve0, _ = reserve0Big.Float64()
	reserve0 /= 1e9

	reserve1Big, err := nextSl.LoadBigCoins()
	if err != nil {
		return nil
	}

	reserve1, _ = reserve1Big.Float64()
	reserve1 /= 1e9

	fmt.Println("[SCN] FOUND DEDUST SWAP")
	fmt.Printf("[SCN] ASSET IN - [%0.2f] [%s]\n", amountIn, assetIn)
	fmt.Printf("[SCN] ASSET OUT - [%0.2f] [%s]\n", amountOut, assetOut)
	fmt.Printf("[SCN] SENDER - [%s]\n", senderAddr)
	fmt.Printf("[SCN] RESERVE 0 - [%0.2f]\n", reserve0)
	fmt.Printf("[SCN] RESERVE 1 - [%0.2f]\n", reserve1)

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
