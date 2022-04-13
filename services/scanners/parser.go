package scanners

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/google"
	"log"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"oasisTracker/smodels"
	"oasisTracker/smodels/container"
	"time"

	beaconAPI "github.com/oasisprotocol/oasis-core/go/beacon/api"
	"github.com/oasisprotocol/oasis-core/go/common/grpc"
	consensusAPI "github.com/oasisprotocol/oasis-core/go/consensus/api"
	grpcCommon "google.golang.org/grpc"
)

const (
	//precision          = 6

	parserBaseTask             = "base"
	parserBalancesSnapshotTask = "balances_snapshot"

	//defaultFlag         ParseFlag = iota
	baseFlag            ParseFlag = 1
	balanceSnapshotFlag           = baseFlag << 1
	//watcherFlag                   = baseFlag | balanceSnapshotFlag
)

type (
	ParseFlag uint32
	Parser    struct {
		ctx context.Context

		dao       dao.ParserDAO
		api       consensusAPI.ClientBackend
		bAPI      beaconAPI.Backend
		conn      *grpcCommon.ClientConn
		baseEpoch beaconAPI.EpochTime

		container *ParseContainer
	}

	ParseContainer struct {
		blocks          *container.BlocksContainer
		blockSignatures *container.BlockSignatureContainer
		txs             *container.TxsContainer
		balances        *container.AccountsContainer
		rewards         *container.RewardsContainer
	}
)

func NewParser(ctx context.Context, cfg conf.Scanner, d dao.ParserDAO) (*Parser, error) {
	credentials := google.NewDefaultCredentials().TransportCredentials()
	grpcConn, err := grpc.Dial(cfg.NodeConfig, grpcCommon.WithTransportCredentials(credentials))
	if err != nil {
		return nil, err
	}

	cAPI := consensusAPI.NewConsensusClient(grpcConn)
	bAPI := beaconAPI.NewBeaconClient(grpcConn)

	baseEpoch, err := bAPI.GetBaseEpoch(ctx)
	if err != nil {
		return nil, err
	}

	return &Parser{
		ctx:       ctx,
		conn:      grpcConn,
		api:       cAPI,
		bAPI:      bAPI,
		dao:       d,
		baseEpoch: baseEpoch,
		container: &ParseContainer{
			blocks:          container.NewBlocksContainer(),
			blockSignatures: container.NewBlockSignatureContainer(),
			txs:             container.NewTxsContainer(),
			balances:        container.NewAccountsContainer(),
			rewards:         container.NewRewardsContainer(),
		},
	}, nil
}

func (p *Parser) GetTaskExecutor(taskTitle string) (executor *smodels.Executor, err error) {
	switch taskTitle {
	case parserBaseTask:
		return &smodels.Executor{
			ExecHeight: p.ParseBase,
			Truncate:   p.Truncate,
			Save:       p.Save,
		}, nil
	case parserBalancesSnapshotTask:
		return &smodels.Executor{
			ExecHeight: p.ParseBalancesSnapshot,
			Truncate:   p.Truncate,
			Save:       p.Save,
		}, nil
	default:
		return nil, fmt.Errorf("executor %s not found", taskTitle)
	}
}

func (p *Parser) Truncate() {

	p.container.blocks.Flush()
	p.container.blockSignatures.Flush()
	p.container.txs.Flush()
	p.container.balances.Flush()
	p.container.rewards.Flush()

}

func (p *Parser) Save() (err error) {
	log.Print("Start saving")
	if !p.container.blocks.IsEmpty() {
		tm := time.Now()
		err := p.dao.CreateBlocks(p.container.blocks.Blocks())
		if err != nil {
			return fmt.Errorf("dao.CreateBlocks: %s", err.Error())
		}

		log.Print("Save time Blocks: ", time.Since(tm))
	}

	if !p.container.blockSignatures.IsEmpty() {
		tm := time.Now()
		err = p.dao.CreateBlockSignatures(p.container.blockSignatures.Signatures())
		if err != nil {
			return fmt.Errorf("dao.CreateBlockSignatures: %s", err.Error())
		}
		log.Print("Save time Signatures: ", time.Since(tm))
	}

	if !p.container.txs.IsEmpty() {
		tm := time.Now()
		err = p.dao.CreateTransfers(p.container.txs.Txs())
		if err != nil {
			return fmt.Errorf("dao.CreateTransfers: %s", err.Error())
		}

		log.Print("Save time Transfers: ", time.Since(tm))

		err = p.dao.CreateRegisterNodeTransactions(p.container.txs.NodeRegistryTxs())
		if err != nil {
			return fmt.Errorf("dao.CreateRegisterNodeTransactions: %s", err.Error())
		}

		err = p.dao.CreateRegisterEntityTransactions(p.container.txs.EntityRegistryTxs())
		if err != nil {
			return fmt.Errorf("dao.CreateRegisterEntityTransactions: %s", err.Error())
		}
	}

	if !p.container.balances.IsEmpty() {
		tm := time.Now()
		err = p.dao.CreateAccountBalances(p.container.balances.Balances())
		if err != nil {
			return fmt.Errorf("dao.CreateAccountBalances: %s", err.Error())
		}

		log.Print("Save time Balances: ", time.Since(tm))
	}

	if !p.container.rewards.IsEmpty() {
		tm := time.Now()
		err = p.dao.CreateRewards(p.container.rewards.Rewards())
		if err != nil {
			return fmt.Errorf("dao.CreateRewards: %s", err.Error())
		}

		log.Print("Save time Rewards: ", time.Since(tm))
	}

	return nil
}

func (p *Parser) ParseBase(conn *grpcCommon.ClientConn, blockID uint64) error {

	parsTask, err := NewParserTask(p.ctx, conn, p.baseEpoch, p.container)
	if err != nil {
		return err
	}

	err = parsTask.ParseBase(blockID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) ParseBalancesSnapshot(conn *grpcCommon.ClientConn, epoch uint64) error {

	parsTask, err := NewParserTask(p.ctx, conn, p.baseEpoch, p.container)
	if err != nil {
		return err
	}

	err = parsTask.EpochBalanceSnapshot(beaconAPI.EpochTime(epoch))
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) ParseWatchBlock(block *consensusAPI.Block) error {
	parsTask, err := NewParserTask(p.ctx, p.conn, p.baseEpoch, p.container)
	if err != nil {
		return err
	}

	err = parsTask.parseOasisBase(block, baseFlag)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) ParseEpochSnap(epoch beaconAPI.EpochTime) error {
	parsTask, err := NewParserTask(p.ctx, p.conn, p.baseEpoch, p.container)
	if err != nil {
		return err
	}

	err = parsTask.EpochBalanceSnapshot(epoch)
	if err != nil {
		return err
	}

	return nil
}
