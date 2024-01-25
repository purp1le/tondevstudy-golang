package scan

import (
	"time"
	"ton-lessons/internal/storage"
	"ton-lessons/internal/structures"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"gorm.io/gorm"
)

func (s *Scanner) processDedustSwapEvent(
	dbtx *gorm.DB,
	transaction *tlb.ExternalMessageOut,
) error {
	var (
		dedustSwap storage.DedustSwap
		assetIn    string
		assetOut   string
	)

	if transaction.Body == nil {
		return nil
	}

	bodySlice := transaction.Body.BeginParse()
	opcode, err := bodySlice.LoadUInt(32)
	if err != nil {
		return nil
	}

	if opcode != 0x9c610de3 {
		return nil
	}

	logrus.Info("FOUND DEDUST SWAP OPCODE")

	isNativeIn, err := bodySlice.LoadUInt(4)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	if isNativeIn == 0 {
		assetIn = "TON"
	} else {
		var (
			worckhain   int64
			addressData []byte
		)

		worckhain, err = bodySlice.LoadInt(8)
		if err != nil {
			logrus.Error(err)
			return nil
		}

		addressData, err = bodySlice.LoadSlice(256)
		if err != nil {
			logrus.Error(err)
			return nil
		}

		assetIn = address.NewAddress(0, byte(worckhain), addressData).String()
	}

	isNativeOut, err := bodySlice.LoadUInt(4)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	if isNativeOut == 0 {
		assetOut = "TON"
	} else {
		var (
			worckhain   int64
			addressData []byte
		)

		worckhain, err = bodySlice.LoadInt(8)
		if err != nil {
			logrus.Error(err)
			return nil
		}

		addressData, err = bodySlice.LoadSlice(256)
		if err != nil {
			logrus.Error(err)
			return nil
		}

		assetOut = address.NewAddress(0, byte(worckhain), addressData).String()
	}

	amountOut, err := bodySlice.LoadBigCoins()
	if err != nil {
		logrus.Error(err)
		return nil
	}

	amountIn, err := bodySlice.LoadBigCoins()
	if err != nil {
		logrus.Error(err)
		return nil
	}

	ref, err := bodySlice.LoadRef()
	if err != nil {
		logrus.Error(err)
		return nil
	}

	senderAddress, err := ref.LoadAddr()
	if err != nil {
		logrus.Error(err)
		return nil
	}

	refAddress, err := ref.LoadAddr()
	if err != nil {
		logrus.Error(err)
		return nil
	}

	reserveLeft, err := ref.LoadBigCoins()
	if err != nil {
		logrus.Error(err)
		return nil
	}

	reserverRight, err := ref.LoadBigCoins()
	if err != nil {
		logrus.Error(err)
		return nil
	}

	dedustSwap = storage.DedustSwap{
		CreatedAt:       time.Unix(int64(transaction.CreatedAt), 0),
		ProcessedAt:     time.Now(),
		AssetIn:         assetIn,
		AssetOut:        assetOut,
		AmountIn:        amountIn.String(),
		AmountOut:       amountOut.String(),
		SenderAddress:   senderAddress.String(),
		ReferralAddress: refAddress.String(),
		ReserveLeft:     reserveLeft.String(),
		ReserveRight:    reserverRight.String(),
		DedustPool:      transaction.SrcAddr.String(),
	}

	if err := dbtx.Create(&dedustSwap).Error; err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Infof("[SCN] found dedust swap from [%s [%s]] to [%s [%s]] sender address [%s] on pool [%s]",
		assetIn,
		amountIn.String(),
		assetOut,
		amountOut.String(),
		senderAddress.String(),
		dedustSwap.DedustPool,
	)

	return nil
}

func (s *Scanner) processDedustDeposit(
	dbtx *gorm.DB,
	transaction *tlb.ExternalMessageOut,
) error {
	var (
		dedustDeposit structures.DedustDeposit
	)

	if transaction.Body == nil {
		return nil
	}

	if err := tlb.LoadFromCell(
		&dedustDeposit,
		transaction.Body.BeginParse(),
		false,
	); err != nil {
		return nil
	}

	deposit := storage.DedustDeposit{
		CreatedAt:     time.Unix(int64(transaction.CreatedAt), 0),
		ProcessedAt:   time.Now(),
		Liquidity:     dedustDeposit.Liquidity.String(),
		AmountIn:      dedustDeposit.AmountLeft.String(),
		AmountOut:     dedustDeposit.AmountRight.String(),
		ReserveLeft:   dedustDeposit.ReserveLeft.String(),
		ReserveRight:  dedustDeposit.ReserveRight.String(),
		SenderAddress: dedustDeposit.SenderAddress.String(),
		DedustPool:    transaction.SrcAddr.String(),
	}

	if err := dbtx.Create(&deposit).Error; err != nil {
		return err
	}

	logrus.Infof("[SCN] found dedust deposit liquidity Amount0 [%s] Amount1 [%s] on pool [%s]",
		dedustDeposit.AmountLeft.String(),
		dedustDeposit.AmountRight.String(),
		transaction.SrcAddr.String(),
	)

	return nil
}

func (s *Scanner) processDedustWithdraw(
	dbtx *gorm.DB,
	transaction *tlb.ExternalMessageOut,
) error {
	var (
		dedustWithdraw structures.DedustWithdraw
	)

	if transaction.Body == nil {
		return nil
	}

	if err := tlb.LoadFromCell(
		&dedustWithdraw,
		transaction.Body.BeginParse(),
		false,
	); err != nil {
		return nil
	}

	deposit := storage.DedustWithdraw{
		CreatedAt:     time.Unix(int64(transaction.CreatedAt), 0),
		ProcessedAt:   time.Now(),
		Liquidity:     dedustWithdraw.Liquidity.String(),
		AmountIn:      dedustWithdraw.AmountLeft.String(),
		AmountOut:     dedustWithdraw.AmountRight.String(),
		ReserveLeft:   dedustWithdraw.ReserveLeft.String(),
		ReserveRight:  dedustWithdraw.ReserveRight.String(),
		SenderAddress: dedustWithdraw.SenderAddress.String(),
		DedustPool:    transaction.SrcAddr.String(),
	}

	if err := dbtx.Create(&deposit).Error; err != nil {
		return err
	}

	logrus.Infof("[SCN] found dedust withdraw liquidity Amount0 [%s] Amount1 [%s] on pool [%s]",
		dedustWithdraw.AmountLeft.String(),
		dedustWithdraw.AmountRight.String(),
		transaction.SrcAddr.String(),
	)

	return nil
}
