package storage

import "time"

type Block struct {
	SeqNo       uint32 `gorm:"primaryKey;autoIncrement:false;"`
	WorkChain   int32
	Shard       int64
	ProcessedAt time.Time
}
