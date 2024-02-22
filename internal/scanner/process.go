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
	switch trans.IO.In.MsgType {
	case tlb.MsgTypeInternal:
		inTrans := trans.IO.In.AsInternal()

		if err := s.processNftContract(
			dbtx,
			master,
			inTrans.DstAddr,
		); err != nil {
			return err
		}
	case tlb.MsgTypeExternalIn:
		extInTrans := trans.IO.In.AsExternalIn()

		if err := s.processNftContract(
			dbtx,
			master,
			extInTrans.DstAddr,
		); err != nil {
			return err
		}
	}

	return nil
}
