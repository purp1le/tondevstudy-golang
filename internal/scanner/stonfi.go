package scan

import (
	"context"
	"encoding/base64"
	"fmt"
	"ton-lessons/internal/app"
	"ton-lessons/internal/structures"

	"github.com/xssnick/tonutils-go/tlb"
)

func (s *Scanner) ProcessStonfiSwapPart1(
	inMsg *tlb.InternalMessage,
	master *tlb.BlockInfo,
) (*structures.StonfiSwapPart1, error) {
	var (
		stonfiSwap1 structures.StonfiSwapPart1
	)

	if inMsg.Body == nil {
		return nil, nil
	}

	if err := tlb.LoadFromCell(
		&stonfiSwap1,
		inMsg.Body.BeginParse(),
		false,
	); err != nil {
		return nil, nil
	}

	accountInfo, err := s.Api.GetAccount(context.Background(), master, inMsg.DstAddr)
	if err != nil {
		return nil, err
	}

	if app.StonfiPoolCodeHash != base64.StdEncoding.EncodeToString(accountInfo.Code.Hash()) {
		return nil, nil
	}

	return &stonfiSwap1, nil
}

func (s *Scanner) ProcessStonfiSwapPart2(
	stonfiSwap1 *structures.StonfiSwapPart1,
	outMsg *tlb.InternalMessage,
) error {
	var (
		stonfiSwap2 structures.StonfiSwapPart2
	)

	if outMsg.Body == nil {
		return nil
	}

	if err := tlb.LoadFromCell(
		&stonfiSwap2,
		outMsg.Body.BeginParse(),
		false,
	); err != nil {
		return nil
	}

	fmt.Println("[SCN] FOUND STONFI SWAP")
	fmt.Println(stonfiSwap1)
	fmt.Println(stonfiSwap2)

	return nil
}
