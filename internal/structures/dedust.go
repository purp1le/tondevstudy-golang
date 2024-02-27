package structures

import (
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type DedustAssetNative struct {
	_ tlb.Magic `tlb:"$0000"`
}

func (a DedustAssetNative) Type() string {
	return "native"
}

func (a DedustAssetNative) AsNative() DedustAssetNative {
	return a
}

func (a DedustAssetNative) AsJetton() DedustAssetJetton {
	return DedustAssetJetton{}
}

type DedustAssetJetton struct {
	_           tlb.Magic `tlb:"$0001"`
	WorkchainId uint64    `tlb:"## 8"`
	AddressData []byte    `tlb:"bits 256"`
}

func (a DedustAssetJetton) Type() string {
	return "jetton"
}

func (a DedustAssetJetton) AsNative() DedustAssetNative {
	return DedustAssetNative{}
}

func (a DedustAssetJetton) AsJetton() DedustAssetJetton {
	return a
}

type DedustAsset interface {
	Type() string
	AsNative() DedustAssetNative
	AsJetton() DedustAssetJetton
}

type DedustSwapEvent struct {
	_         tlb.Magic   `tlb:"#9c610de3"`
	AssetIn   DedustAsset `tlb:"[DedustAssetJetton,DedustAssetNative]"`
	AssetOut  DedustAsset `tlb:"[DedustAssetJetton,DedustAssetNative]"`
	AmountIn  *tlb.Coins  `tlb:"."`
	AmountOut *tlb.Coins  `tlb:"."`
	ExtraInfo struct {
		SenderAddr   *address.Address `tlb:"addr"`
		ReferralAddr *address.Address `tlb:"addr"`
		Reserve0     *tlb.Coins       `tlb:"."`
		Reserve1     *tlb.Coins       `tlb:"."`
	} `tlb:"^"`
}

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

type DedustSwapStepParams struct {
	_     tlb.Magic       `tlb:"$0"`
	Limit tlb.Coins       `tlb:"."`
	Next  *DedustSwapStep `tlb:"maybe ^"`
}

type DedustSwapStep struct {
	PoolAddr       *address.Address     `tlb:"addr"`
	SwapStepParams DedustSwapStepParams `tlb:"."`
}

type DedustSwapParams struct {
	Deadline       uint32           `tlb:"## 32"`
	RecipientAddr  *address.Address `tlb:"addr"`
	ReferralAddr   *address.Address `tlb:"addr"`
	FulfillPayload *cell.Cell       `tlb:"maybe ^"`
	RejectPayload  *cell.Cell       `tlb:"maybe ^"`
}

type DedustRequestNativeSwap struct {
	_          tlb.Magic        `tlb:"#ea06185d"`
	QueryId    uint64           `tlb:"## 64"`
	Amount     tlb.Coins        `tlb:"."`
	SwapStep   DedustSwapStep   `tlb:"."`
	SwapParams DedustSwapParams `tlb:"^"`
}

type DedustRequestJettonSwap struct {
	_          tlb.Magic        `tlb:"#e3a0d482"`
	SwapStep   DedustSwapStep   `tlb:"."`
	SwapParams DedustSwapParams `tlb:"^"`
}
