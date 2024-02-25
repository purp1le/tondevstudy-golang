package structures

import (
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
)

type StonfiSwapPart1 struct {
	_            tlb.Magic        `tlb:"#25938561"`
	QueryId      uint64           `tlb:"## 64"`
	TokenWallet1 *address.Address `tlb:"addr"`
	MinOut       *tlb.Coins       `tlb:"."`
	ToAddress    *address.Address `tlb:"addr"`
	HasRef       bool             `tlb:"bool"`
	RefAddress   *address.Address `tlb:"?HasRef addr"`
}

type StonfiSwapPart2 struct {
	_            tlb.Magic        `tlb:"#f93bb43f"`
	QueryId      uint64           `tlb:"## 64"`
	OwnerAddr    *address.Address `tlb:"addr"`
	ExitCode     uint32           `tlb:"## 32"`
	RefCoinsData struct {
		Amount0Out    *tlb.Coins       `tlb:"."`
		Token0Address *address.Address `tlb:"addr"`
		Amount1Out    *tlb.Coins       `tlb:"."`
		Token1Address *address.Address `tlb:"addr"`
	} `tlb:"^"`
}
