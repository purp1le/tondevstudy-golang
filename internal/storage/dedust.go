package storage

import "time"

type DedustSwap struct {
	Id              uint64 `gorm:"primaryKey;autoIncrement:true;"`
	CreatedAt       time.Time
	ProcessedAt     time.Time
	AssetIn         string // ton | jetton address
	AssetOut        string // ton | jetton address
	AmountIn        string
	AmountOut       string
	SenderAddress   string
	ReferralAddress string
	ReserveLeft     string
	ReserveRight    string
	DedustPool      string
}

type DedustDeposit struct {
	Id            uint64 `gorm:"primaryKey;autoIncrement:true;"`
	CreatedAt     time.Time
	ProcessedAt   time.Time
	Liquidity     string
	AmountIn      string
	AmountOut     string
	ReserveLeft   string
	ReserveRight  string
	SenderAddress string
	DedustPool    string
}

type DedustWithdraw struct {
	Id            uint64 `gorm:"primaryKey;autoIncrement:true;"`
	CreatedAt     time.Time
	ProcessedAt   time.Time
	Liquidity     string
	AmountIn      string
	AmountOut     string
	ReserveLeft   string
	ReserveRight  string
	SenderAddress string
	DedustPool    string
}
