package app

import (
	"ton-lessons/internal/structures"

	"github.com/xssnick/tonutils-go/tlb"
)

func InitTlb() {
	tlb.Register(structures.DedustAssetNative{})
	tlb.Register(structures.DedustAssetJetton{})
}
