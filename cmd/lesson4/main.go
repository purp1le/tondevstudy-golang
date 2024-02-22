package main

import (
	"context"
	"encoding/hex"
	"math/big"
	"ton-lessons/internal/app"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/nft"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := app.InitApp(); err != nil {
		return err
	}

	liteclient := liteclient.NewConnectionPool()
	if err := liteclient.AddConnectionsFromConfig(context.Background(), app.CFG.MAINNET_CONFIG); err != nil {
		logrus.Error(err)
		return err
	}

	api := ton.NewAPIClient(liteclient)

	seed := app.CFG.Wallet.SEED
	wall, err := wallet.FromSeed(api, seed, wallet.HighloadV2Verified)
	if err != nil {
		return err
	}

	collectionAddress := address.MustParseAddr("EQAzErrrzcWcD5xgOPTjRW8fN4G8C4J8C67d6WQL6WklchC_")

	collectionClient := nft.NewCollectionClient(api, collectionAddress)

	nftMintPayload, err := collectionClient.BuildMintPayload(big.NewInt(0), wall.Address(), tlb.MustFromTON("0.1"), &nft.ContentOnchain{
		Name:        "new item",
		Description: "example description",
		Image:       "https://avatars.mds.yandex.net/i?id=fd233f798f5a49a1c8bca7fc16db66d6c758d918-10807817-images-thumbs&n=13",
	})

	if err != nil {
		return err
	}

	if err := wall.Send(
		context.Background(),
		wallet.SimpleMessage(
			collectionAddress,
			tlb.MustFromTON("0.2"),
			nftMintPayload,
		), true,
	); err != nil {
		return err
	}

	return nil
}

func getNFTCollectionCode() *cell.Cell {
	var hexBOC = "b5ee9c72410213010001fe000114ff00f4a413f4bcf2c80b0102016204020201200e030025bc82df6a2687d20699fea6a6a182de86a182c40202cd0a050201200706003d45af0047021f005778018c8cb0558cf165004fa0213cb6b12ccccc971fb0080201200908001b3e401d3232c084b281f2fff27420002d007232cffe0a33c5b25c083232c044fd003d0032c0326003ebd10638048adf000e8698180b8d848adf07d201800e98fe99ff6a2687d20699fea6a6a184108349e9ca829405d47141baf8280e8410854658056b84008646582a802e78b127d010a65b509e58fe59f80e78b64c0207d80701b28b9e382f970c892e000f18112e001718119026001f1812f82c207f97840d0c0b002801fa40304144c85005cf1613cb3fccccccc9ed5400a6357003d4308e378040f4966fa5208e2906a4208100fabe93f2c18fde81019321a05325bbf2f402fa00d43022544b30f00623ba9302a402de04926c21e2b3e6303250444313c85005cf1613cb3fccccccc9ed5400603502d33f5313bbf2e1925313ba01fa00d43028103459f0068e1201a44343c85005cf1613cb3fccccccc9ed54925f05e2020120120f0201201110002db4f47da89a1f481a67fa9a9a86028be09e008e003e00b0002fb5dafda89a1f481a67fa9a9a860d883a1a61fa61ff4806100043b8b5d31ed44d0fa40d33fd4d4d43010245f04d0d431d430d071c8cb0701cf16ccc98f34ea10e"
	codeCellBytes, _ := hex.DecodeString(hexBOC)

	codeCell, err := cell.FromBOC(codeCellBytes)
	if err != nil {
		panic(err)
	}

	return codeCell
}

func getNFTItemCode() *cell.Cell {
	var hexBOC = "b5ee9c7241020d010001d0000114ff00f4a413f4bcf2c80b0102016203020009a11f9fe0050202ce050402012008060201200907001d00f232cfd633c58073c5b3327b552000113e910c1c2ebcb85360003b3b513434cffe900835d27080269fc07e90350c04090408f80c1c165b5b6002d70c8871c02497c0f83434c0c05c6c2497c0f83e903e900c7e800c5c75c87e800c7e800c3c00812ce3850c1b088d148cb1c17cb865407e90350c0408fc00f801b4c7f4cfe08417f30f45148c2ea3a1cc840dd78c9004f80c0d0d0d4d60840bf2c9a884aeb8c097c12103fcbc200b0a00727082108b77173505c8cbff5004cf1610248040708010c8cb055007cf165005fa0215cb6a12cb1fcb3f226eb39458cf17019132e201c901fb0001f65135c705f2e191fa4021f001fa40d20031fa00820afaf0801ba121945315a0a1de22d70b01c300209206a19136e220c2fff2e192218e3e821005138d91c85009cf16500bcf16712449145446a0708010c8cb055007cf165005fa0215cb6a12cb1fcb3f226eb39458cf17019132e201c901fb00104794102a375be20c0082028e3526f0018210d53276db103744006d71708010c8cb055007cf165005fa0215cb6a12cb1fcb3f226eb39458cf17019132e201c901fb0093303234e25502f003cc82807e"
	codeCellBytes, _ := hex.DecodeString(hexBOC)

	codeCell, err := cell.FromBOC(codeCellBytes)
	if err != nil {
		panic(err)
	}

	return codeCell
}

func getContractData(collectionOwnerAddr, royaltyAddr *address.Address) *cell.Cell {
	// storage schema
	// default#_ royalty_factor:uint16 royalty_base:uint16 royalty_address:MsgAddress = RoyaltyParams;
	// storage#_ owner_address:MsgAddress next_item_index:uint64
	//           ^[collection_content:^Cell common_content:^Cell]
	//           nft_item_code:^Cell
	//           royalty_params:^RoyaltyParams
	//           = Storage;

	royalty := cell.BeginCell().
		MustStoreUInt(5, 16). // 5% royalty
		MustStoreUInt(100, 16).
		MustStoreAddr(royaltyAddr).
		EndCell()

	// collection data
	collectionContent := nft.ContentOnchain{
		Name:        "Ton lessons collection",
		Description: "Example ton lesson collection",
	}
	collectionContentCell, err := collectionContent.ContentCell()
	if err != nil {
		logrus.Error(err)
		panic(err)
	}

	// prefix for NFTs data
	uri := ""
	commonContentCell := cell.BeginCell().MustStoreStringSnake(uri).EndCell()

	contentRef := cell.BeginCell().
		MustStoreRef(collectionContentCell).
		MustStoreRef(commonContentCell).
		EndCell()

	data := cell.BeginCell().MustStoreAddr(collectionOwnerAddr).
		MustStoreUInt(0, 64).
		MustStoreRef(contentRef).
		MustStoreRef(getNFTItemCode()).
		MustStoreRef(royalty).
		EndCell()

	return data
}
