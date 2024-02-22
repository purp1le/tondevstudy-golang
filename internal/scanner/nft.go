package scan

import (
	"context"
	"time"
	"ton-lessons/internal/storage"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/nft"
	"gorm.io/gorm"
)

func (s *Scanner) processNftContract(
	dbtx *gorm.DB,
	master *tlb.BlockInfo,
	dstAddr *address.Address,
) error {

	if err := s.processNftMaster(
		dbtx,
		master,
		dstAddr,
	); err != nil {
		return err
	}

	if err := s.processNftItem(
		dbtx,
		master,
		dstAddr,
	); err != nil {
		return err
	}

	return nil
}

func (s *Scanner) processNftMaster(
	dbtx *gorm.DB,
	master *tlb.BlockInfo,
	dstAddr *address.Address,
) error {
	var (
		collectionDB storage.NftCollection
	)

	collection := nft.NewCollectionClient(s.Api, dstAddr)

	collectionData, err := collection.GetCollectionDataAtBlock(context.Background(), master)
	if err != nil {
		return nil
	}

	err = dbtx.Where("address = ?", dstAddr.String()).First(&collectionDB).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == nil {
		return nil
	}

	collectionDB = storage.NftCollection{
		Address:        dstAddr.String(),
		Owner:          collectionData.OwnerAddress.String(),
		MintedItemsBig: collectionData.NextItemIndex.Bytes(),
		FoundAt:        time.Now(),
	}

	if err := dbtx.Create(&collectionDB).Error; err != nil {
		return err
	}

	logrus.Infof("[SCN] found new collection not in db [%s]", dstAddr.String())

	return nil
}

func (s *Scanner) processNftItem(
	dbtx *gorm.DB,
	master *tlb.BlockInfo,
	dstAddr *address.Address,
) error {
	var (
		nftItemDB       storage.NftItem
		nftCollectionDB storage.NftCollection
	)

	nftClient := nft.NewItemClient(s.Api, dstAddr)

	nftItem, err := nftClient.GetNFTDataAtBlock(context.Background(), master)
	if err != nil {
		return nil
	}

	err = dbtx.Where("address = ?", dstAddr.String()).First(&nftItemDB).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == nil {
		if nftItemDB.Owner == nftItem.OwnerAddress.String() {
			return nil
		}

		nftItemDB.Owner = nftItem.OwnerAddress.String()
		if err := dbtx.Save(&nftItemDB).Error; err != nil {
			return err
		}

		logrus.Infof("[SCN] nft item change owner [%s]", nftItemDB.Address)

		return nil
	}

	err = dbtx.Where("address = ?", nftItem.CollectionAddress.String()).First(&nftCollectionDB).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		if err := s.processNftMaster(
			dbtx,
			master,
			nftItem.CollectionAddress,
		); err != nil {
			return err
		}

		err = dbtx.Where("address = ?", nftItem.CollectionAddress.String()).First(&nftCollectionDB).Error
		if err != nil {
			return err
		}
	}

	nftItemDB = storage.NftItem{
		CollectionId: nftCollectionDB.Id,
		IndexBig:     nftItem.Index.Bytes(),
		Address:      dstAddr.String(),
		Owner:        nftItem.OwnerAddress.String(),
	}

	if err := dbtx.Create(&nftItemDB).Error; err != nil {
		return err
	}

	logrus.Infof("[SCN] add new item [%s]", nftItemDB.Address)

	return nil
}
