package scan

import (
	"context"
	"encoding/base64"
	"ton-lessons/internal/app"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/tlb"
	"gorm.io/gorm"
)

func (s *Scanner) processTransaction(dbtx *gorm.DB, trans *tlb.Transaction, master *tlb.BlockInfo) error {
	if trans.IO.Out == nil {
		return nil
	}

	messages, err := trans.IO.Out.ToSlice()
	if err != nil {
		return nil
	}

	for _, msg := range messages {
		if msg.MsgType != tlb.MsgTypeExternalOut {
			continue
		}

		msgOut := msg.AsExternalOut()

		accountInfo, err := s.Api.GetAccount(context.Background(), master, msgOut.SrcAddr)
		if err != nil {
			return err
		}

		codeHash := base64.StdEncoding.EncodeToString(accountInfo.Code.Hash())
		if codeHash != app.DedustPoolCodeHash {
			logrus.Info("invalid code")
			continue
		}

		logrus.Info("[SCN] contract code hash: ", codeHash)
		logrus.Info("[SCN] init code hash: ", app.DedustPoolCodeHash)
		if err := s.processDedustSwapEvent(
			dbtx,
			msgOut,
		); err != nil {
			return err
		}

		if err := s.processDedustDeposit(
			dbtx,
			msgOut,
		); err != nil {
			return err
		}

		if err := s.processDedustWithdraw(
			dbtx,
			msgOut,
		); err != nil {
			return err
		}
	}

	return nil
}

// func (s *Scanner) processJettonTransferRequest(
// 	dbtx *gorm.DB,
// 	trans *tlb.InternalMessage,
// ) error {
// 	var (
// 		transferRequest structures.TransferRequest
// 	)

// 	if trans.Body == nil {
// 		return nil
// 	}

// 	if err := tlb.LoadFromCell(
// 		&transferRequest,
// 		trans.Body.BeginParse(),
// 		false,
// 	); err != nil {
// 		return nil
// 	}

// 	logrus.Infof("[SCN] New transfer request from [%s] to [%s] amount [%s]",
// 		trans.SrcAddr.String(),
// 		transferRequest.Destination,
// 		transferRequest.Amount,
// 	)
// 	return nil
// }

// func (s *Scanner) processJettonTransferNotifications(
// 	dbtx *gorm.DB,
// 	trans *tlb.InternalMessage,
// ) error {
// 	var (
// 		transferNotification structures.TransferNotification
// 	)

// 	if trans.Body == nil {
// 		return nil
// 	}

// 	if err := tlb.LoadFromCell(
// 		&transferNotification,
// 		trans.Body.BeginParse(),
// 		false,
// 	); err != nil {
// 		return nil
// 	}

// 	logrus.Infof("[SCN] New transfer notification from [%s] to [%s] amount [%s]",
// 		transferNotification.Sender,
// 		trans.DstAddr.String(),
// 		transferNotification.Amount,
// 	)

// 	return nil
// }
