package structures

import (
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
)

type DedustDeposit struct {
	_             tlb.Magic        `tlb:"#b544f4a4"`
	SenderAddress *address.Address `tlb:"addr"`
	AmountLeft    *tlb.Coins       `tlb:"."`
	AmountRight   *tlb.Coins       `tlb:"."`
	ReserveLeft   *tlb.Coins       `tlb:"."`
	ReserveRight  *tlb.Coins       `tlb:"."`
	Liquidity     *tlb.Coins       `tlb:"."`
}

type DedustWithdraw struct {
	_             tlb.Magic        `tlb:"#3aa870a6"`
	SenderAddress *address.Address `tlb:"addr"`
	Liquidity     *tlb.Coins       `tlb:"."`
	AmountLeft    *tlb.Coins       `tlb:"."`
	AmountRight   *tlb.Coins       `tlb:"."`
	ReserveLeft   *tlb.Coins       `tlb:"."`
	ReserveRight  *tlb.Coins       `tlb:"."`
}
