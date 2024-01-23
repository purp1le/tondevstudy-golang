package scan

import (
	"context"
	"fmt"
	"time"
	"ton-lessons/internal/storage"

	"github.com/xssnick/tonutils-go/ton"
	"gorm.io/gorm"
)

func getShardID(shard *ton.BlockIDExt) string {
	return fmt.Sprintf("%d|%d", shard.Workchain, shard.Shard)
}

func (s *Scanner) getNotSeenShards(ctx context.Context, api *ton.APIClient, shard *ton.BlockIDExt, shardLastSeqno map[string]uint32) (ret []*ton.BlockIDExt, err error) {
	if no, ok := shardLastSeqno[getShardID(shard)]; ok && no == shard.SeqNo {
		return nil, nil
	}

	b, err := api.GetBlockData(ctx, shard)
	if err != nil {
		return nil, fmt.Errorf("get block data: %w", err)
	}

	parents, err := b.BlockInfo.GetParentBlocks()
	if err != nil {
		return nil, fmt.Errorf("get parent blocks (%d:%x:%d): %w", shard.Workchain, uint64(shard.Shard), shard.Shard, err)
	}

	for _, parent := range parents {
		ext, err := s.getNotSeenShards(ctx, api, parent, shardLastSeqno)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ext...)
	}

	ret = append(ret, shard)
	return ret, nil
}

func (s *Scanner) getLastBlockSeqno() (uint32, error) {
	lastMaster, err := s.Api.GetMasterchainInfo(context.Background())
	if err != nil {
		return 0, err
	}

	return lastMaster.SeqNo, nil
}

func (s *Scanner) addBlock(master ton.BlockIDExt, dbtx *gorm.DB) error {
	newBlock := storage.Block{
		SeqNo:       master.SeqNo,
		WorkChain:   master.Workchain,
		Shard:       master.Shard,
		ProcessedAt: time.Now(),
	}

	if err := dbtx.Create(&newBlock).Error; err != nil {
		return err
	}

	s.LastBlock = newBlock
	s.LastBlock.SeqNo += 1
	return nil
}
