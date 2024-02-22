package storage

import "time"

type NftCollection struct {
	Id             uint64
	Address        string
	Owner          string
	MintedItemsBig []byte
	FoundAt        time.Time
}

type NftItem struct {
	Id           uint64
	CollectionId uint64
	IndexBig     []byte
	Address      string
	Owner        string
	FoundAt      time.Time
}
