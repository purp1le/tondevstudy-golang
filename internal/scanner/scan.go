package scan

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
	"ton-lessons/internal/app"
	"ton-lessons/internal/storage"

	"github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"gopkg.in/tomb.v1"
)

type Scanner struct {
	Api            *ton.APIClient
	LastBlock      storage.Block
	Log            *logrus.Logger
	shardLastSeqno map[string]uint32
}

func NewScanner() (*Scanner, error) {
	client := liteclient.NewConnectionPool()

	if err := client.AddConnectionsFromConfig(context.Background(), app.CFG.MAINNET_CONFIG); err != nil {
		return nil, err
	}

	api := ton.NewAPIClient(client)

	var lastBlock storage.Block
	if err := app.DB.Last(&lastBlock).Error; err != nil {
		logrus.Debug("not found blocks in DB")
	}

	if lastBlock.SeqNo != 0 {
		lastBlock.SeqNo += 1
	}

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	logger.SetLevel(logrus.DebugLevel)

	txtFormatter := &logrus.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		ForceColors:            true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", app.FormatFilePath(f.File), f.Line)
		},
	}

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(txtFormatter)
	logger.SetReportCaller(true)

	return &Scanner{
		Api:            api,
		LastBlock:      lastBlock,
		Log:            logger,
		shardLastSeqno: make(map[string]uint32),
	}, nil
}

func (s *Scanner) Listen() {
	s.Log.Debug("[SCN] start processing blocks")
	lastMaster, err := s.Api.GetMasterchainInfo(context.Background())
	for err != nil {
		lastMaster, err = s.Api.GetMasterchainInfo(context.Background())
		s.Log.Error(err)
	}

	if s.LastBlock.SeqNo == 0 && app.CFG.START_BLOCK != 0 {
	} else if s.LastBlock.SeqNo == 0 {
		s.LastBlock.SeqNo = lastMaster.SeqNo
	}

	s.LastBlock.Shard = lastMaster.Shard
	s.LastBlock.WorkChain = lastMaster.Workchain

	master, err := s.Api.LookupBlock(context.Background(), s.LastBlock.WorkChain, s.LastBlock.Shard, s.LastBlock.SeqNo)
	for err != nil {
		s.Log.Error(err)
		time.Sleep(time.Second * 2)
		master, err = s.Api.LookupBlock(context.Background(), s.LastBlock.WorkChain, s.LastBlock.Shard, s.LastBlock.SeqNo)
	}

	firstShards, err := s.Api.GetBlockShardsInfo(context.Background(), master)
	for err != nil {
		s.Log.Error(err)
		time.Sleep(time.Second * 2)
		firstShards, err = s.Api.GetBlockShardsInfo(context.Background(), master)
	}

	for _, shard := range firstShards {
		s.shardLastSeqno[getShardID(shard)] = shard.SeqNo
	}

	s.processBlocks()
}

func (s *Scanner) processBlocks() {
	for {
		master, err := s.Api.LookupBlock(context.Background(), s.LastBlock.WorkChain, s.LastBlock.Shard, s.LastBlock.SeqNo)
		for err != nil {
			s.Log.Error("[SCN] lookup block err, sleep 2 sec: ", err)
			time.Sleep(time.Second * 2)
			master, err = s.Api.LookupBlock(context.Background(), s.LastBlock.WorkChain, s.LastBlock.Shard, s.LastBlock.SeqNo)
		}

		scanErr := s.processMcBlock(master)
		for scanErr != nil {
			s.Log.Error("[SCN] mc block process err: ", scanErr)
			time.Sleep(time.Second * 2)
			scanErr = s.processMcBlock(master)
		}
	}
}

func (s *Scanner) processMcBlock(master *ton.BlockIDExt) error {
	timeStart := time.Now()

	// getting information about other work-chains and shards of master block
	currentShards, err := s.Api.GetBlockShardsInfo(context.Background(), master)
	if err != nil {
		return err
	}

	if len(currentShards) == 0 {
		s.Log.Debugf("block [%d] without shards", master.SeqNo)
		return nil
	}

	var newShards []*ton.BlockIDExt

	for _, shard := range currentShards {
		notSeen, err := s.getNotSeenShards(context.Background(), s.Api, shard, s.shardLastSeqno)
		if err != nil {
			return err
		}
		s.shardLastSeqno[getShardID(shard)] = shard.SeqNo
		newShards = append(newShards, notSeen...)
	}

	// case when we have error and try to process this block again, newShards are empty
	if len(newShards) == 0 {
		newShards = currentShards
	}
	newShards = append(newShards, currentShards...)

	var txList []*tlb.Transaction

	// for each shard block getting transactions
	var wg sync.WaitGroup
	var tomb tomb.Tomb
	allDone := make(chan struct{})
	for _, shard := range newShards {

		var fetchedIDs []ton.TransactionShortInfo
		var after *ton.TransactionID3
		var more = true

		// load all transactions in batches with 100 transactions in each while exists
		for more {

			fetchedIDs, more, err = s.Api.GetBlockTransactionsV2(context.Background(), shard, 100, after)
			if err != nil {
				return err
			}

			if more {
				after = fetchedIDs[len(fetchedIDs)-1].ID3()
			}

			for _, id := range fetchedIDs {
				// get full transaction by id
				wg.Add(1)
				go func(shard *tlb.BlockInfo, account []byte, lt uint64) {
					defer wg.Done()
					tx, err := s.Api.GetTransaction(context.Background(), shard, address.NewAddress(0, 0, account), lt)
					if err != nil {
						tomb.Kill(err)
					}
					txList = append(txList, tx)
				}(shard, id.Account, id.LT)
			}
		}
	}

	go func() {
		wg.Wait()
		close(allDone)
	}()

	select {
	case <-allDone:
	case <-tomb.Dying():
		logrus.Error("[SCN] err when get transactions: ", tomb.Err())
		return tomb.Err()
	}
	tomb.Done()

	dbtx := app.DB.Begin()
	for _, transaction := range txList {
		if err := s.processTransaction(dbtx, transaction, master); err != nil {
			logrus.Error(err)
		}
	}

	if err := s.addBlock(*master, dbtx); err != nil {
		dbtx.Rollback()
		return err
	} else {
		if err := dbtx.Commit().Error; err != nil {
			return err
		}
	}

	lastSeqno, err := s.getLastBlockSeqno()
	if err != nil {
		s.Log.Infof("[SCN] success process block [%d] time to process block [%0.2fs] trans count [%d]",
			master.SeqNo,
			time.Since(timeStart).Seconds(),
			len(txList),
		)
	} else {
		s.Log.Infof("[SCN] success process block [%d|%d] time to process block [%0.2fs] trans count [%d]",
			master.SeqNo,
			lastSeqno,
			time.Since(timeStart).Seconds(),
			len(txList),
		)
	}

	return nil
}
