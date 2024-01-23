package scan

import (
	"ton-lessons/internal/structures"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/tlb"
	"gorm.io/gorm"
)

func (s *Scanner) processTransaction(dbtx *gorm.DB, trans *tlb.Transaction) error {

	if trans.IO.In.MsgType == tlb.MsgTypeInternal {
		if err := s.processJettonTransferNotifications(
			dbtx,
			trans.IO.In.AsInternal(),
		); err != nil {
			return err
		}

		if err := s.processJettonTransferRequest(
			dbtx,
			trans.IO.In.AsInternal(),
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scanner) processJettonTransferRequest(
	dbtx *gorm.DB,
	trans *tlb.InternalMessage,
) error {
	var (
		transferRequest structures.TransferRequest
	)

	if trans.Body == nil {
		return nil
	}

	if err := tlb.LoadFromCell(
		&transferRequest,
		trans.Body.BeginParse(),
		false,
	); err != nil {
		return nil
	}

	logrus.Infof("[SCN] New transfer request from [%s] to [%s] amount [%s]",
		trans.SrcAddr.String(),
		transferRequest.Destination,
		transferRequest.Amount,
	)
	return nil
}

func (s *Scanner) processJettonTransferNotifications(
	dbtx *gorm.DB,
	trans *tlb.InternalMessage,
) error {
	var (
		transferNotification structures.TransferNotification
	)

	if trans.Body == nil {
		return nil
	}

	if err := tlb.LoadFromCell(
		&transferNotification,
		trans.Body.BeginParse(),
		false,
	); err != nil {
		return nil
	}

	logrus.Infof("[SCN] New transfer notification from [%s] to [%s] amount [%s]",
		transferNotification.Sender,
		trans.DstAddr.String(),
		transferNotification.Amount,
	)

	return nil
}
