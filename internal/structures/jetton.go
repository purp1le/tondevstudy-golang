package structures

import (
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type JettonTrasfer struct {
	_                   tlb.Magic        `tlb:"#0f8a7ea5"`
	QueryId             uint64           `tlb:"## 64"`
	Amount              tlb.Coins        `tlb:"."`
	Destination         *address.Address `tlb:"addr"`
	ResponseDestination *address.Address `tlb:"addr"`
	CustomPayload       *cell.Cell       `tlb:"maybe ^"`
	FwdTonAmount        tlb.Coins        `tlb:"."`
	FwdPayload          *cell.Cell       `tlb:"either . ^"`
}

type JettonNotification struct {
	_          tlb.Magic        `tlb:"#7362d09c"`
	QueryId    uint64           `tlb:"## 64"`
	Amount     tlb.Coins        `tlb:"."`
	Sender     *address.Address `tlb:"addr"`
	FwdPayload *cell.Cell       `tlb:"either . ^"`
}
