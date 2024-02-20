package structures

import (
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
)

type DedustDepositEvent struct {
	_          tlb.Magic        `tlb:"#b544f4a4"`
	SenderAddr *address.Address `tlb:"addr"`
	Amount0    *tlb.Coins       `tlb:"."`
	Amount1    *tlb.Coins       `tlb:"."`
	Reserve0   *tlb.Coins       `tlb:"."`
	Reserve1   *tlb.Coins       `tlb:"."`
	Liquidity  *tlb.Coins       `tlb:"."`
}

type DedustDepositWithdrawal struct {
	_          tlb.Magic        `tlb:"#3aa870a6"`
	SenderAddr *address.Address `tlb:"addr"`
	Liquidity  *tlb.Coins       `tlb:"."`
	Amount0    *tlb.Coins       `tlb:"."`
	Amount1    *tlb.Coins       `tlb:"."`
	Reserve0   *tlb.Coins       `tlb:"."`
	Reserve1   *tlb.Coins       `tlb:"."`
}
