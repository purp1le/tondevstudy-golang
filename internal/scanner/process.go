package scan

import (
	"ton-lessons/internal/structures"

	"github.com/xssnick/tonutils-go/tlb"
	"gorm.io/gorm"
)

func (s *Scanner) processTransaction(
	dbtx *gorm.DB,
	trans *tlb.Transaction,
	master *tlb.BlockInfo,
) error {
	var (
		stonfiSwap1 *structures.StonfiSwapPart1
	)

	if trans.IO.In.MsgType != tlb.MsgTypeInternal {
		return nil
	}
	stonfiSwap1, err := s.ProcessStonfiSwapPart1(
		trans.IO.In.AsInternal(),
		master,
	)
	if err != nil {
		return err
	}

	if trans.IO.Out == nil {
		return nil
	}

	outMsgs, err := trans.IO.Out.ToSlice()
	if err != nil {
		return err
	}

	for _, msg := range outMsgs {
		switch msg.MsgType {
		case tlb.MsgTypeInternal:
			if stonfiSwap1 != nil {
				if err := s.ProcessStonfiSwapPart2(
					stonfiSwap1,
					msg.AsInternal(),
				); err != nil {
					return err
				}
			}
		case tlb.MsgTypeExternalOut:
			if err := s.processDedustSwap(master, msg.AsExternalOut()); err != nil {
				return err
			}
		}
	}

	// switch trans.IO.In.MsgType {
	// case tlb.MsgTypeInternal:
	// 	inTrans := trans.IO.In.AsInternal()

	// 	if err := s.processNftContract(
	// 		dbtx,
	// 		master,
	// 		inTrans.DstAddr,
	// 	); err != nil {
	// 		return err
	// 	}
	// case tlb.MsgTypeExternalIn:
	// 	extInTrans := trans.IO.In.AsExternalIn()

	// 	if err := s.processNftContract(
	// 		dbtx,
	// 		master,
	// 		extInTrans.DstAddr,
	// 	); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
