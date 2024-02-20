package scan

import (
	"github.com/xssnick/tonutils-go/tlb"
	"gorm.io/gorm"
)

func (s *Scanner) processTransaction(
	dbtx *gorm.DB,
	trans *tlb.Transaction,
	master *tlb.BlockInfo,
) error {
	if trans.IO.Out == nil {
		return nil
	}

	outMsgs, err := trans.IO.Out.ToSlice()
	if err != nil {
		return nil
	}

	for _, msg := range outMsgs {
		if msg.MsgType != tlb.MsgTypeExternalOut {
			continue
		}

		extOutMsg := msg.AsExternalOut()

		if err := s.processDedustSwap(master, extOutMsg); err != nil {
			return err
		}

		if err := s.processDedustDeposit(master, extOutMsg); err != nil {
			return err
		}

		if err := s.processDedustWithdrawal(master, extOutMsg); err != nil {
			return err
		}
	}

	return nil
}

// func (s *Scanner) findJettonTransferRequest(
// 	inTrans *tlb.InternalMessage,
// ) error {
// 	var (
// 		jettonTransferRequest structures.JettonTrasfer
// 	)

// 	if err := tlb.LoadFromCell(
// 		&jettonTransferRequest,
// 		inTrans.Body.BeginParse(),
// 		false,
// 	); err != nil {
// 		return nil
// 	}

// 	logrus.Infof("[SCN] trnsfr rqst from [%s] to [%s] amount [%s]",
// 		inTrans.SrcAddr,
// 		jettonTransferRequest.Destination,
// 		jettonTransferRequest.Amount.String(),
// 	)

// 	return nil
// }

// func (s *Scanner) findJettonTransferNotification(
// 	inTrans *tlb.InternalMessage,
// ) error {
// 	var (
// 		jettonTransferNotification structures.JettonNotification
// 	)

// 	if err := tlb.LoadFromCell(
// 		&jettonTransferNotification,
// 		inTrans.Body.BeginParse(),
// 		false,
// 	); err != nil {
// 		return nil
// 	}

// 	logrus.Infof("[SCN] trnsfr ntfctn from [%s] to [%s] amount [%s]",
// 		jettonTransferNotification.Sender,
// 		inTrans.DstAddr,
// 		jettonTransferNotification.Amount.String(),
// 	)

// 	return nil
// }
