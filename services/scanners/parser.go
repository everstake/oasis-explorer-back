package scanners

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/oasisprotocol/oasis-core/go/common/grpc"
	consensusAPI "github.com/oasisprotocol/oasis-core/go/consensus/api"
	"golang.org/x/crypto/blake2b"
	grpcCommon "google.golang.org/grpc"
	"log"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"oasisTracker/smodels"
	"strconv"
	"time"
)

const (
	precision          = 6
	saveBatch          = 200
	saveAddressesBatch = 50

	parserBaseTask             = "base"
	parserBalancesSnapshotTask = "balances_snapshot"
	parserSignaturesTask       = "signatures"
	parseTransactionsTask      = "transactions"

	defaultFlag         ParseFlag = iota
	baseFlag            ParseFlag = 1
	balanceSnapshotFlag           = baseFlag << 1
	watcherFlag                   = baseFlag | balanceSnapshotFlag
)

type (
	ParseFlag uint32
	Parser    struct {
		ctx context.Context

		dao  dao.ParserDAO
		api  consensusAPI.ClientBackend
		conn *grpcCommon.ClientConn

		container *ParseContainer
	}

	ParseContainer struct {
		blocks          *smodels.BlocksContainer
		blockSignatures *smodels.BlockSignatureContainer
		txs             *smodels.TxsContainer
		balances        *smodels.AccountsContainer
	}
)

func NewParser(ctx context.Context, cfg conf.Scanner, d dao.ParserDAO) (*Parser, error) {
	grpcConn, err := grpc.Dial(cfg.NodeConfig, grpcCommon.WithInsecure())
	if err != nil {
		return nil, err
	}

	cAPI := consensusAPI.NewConsensusClient(grpcConn)

	return &Parser{
		ctx:  ctx,
		conn: grpcConn,
		api:  cAPI,
		dao:  d,
		container: &ParseContainer{
			blocks:          smodels.NewBlocksContainer(),
			blockSignatures: smodels.NewBlockSignatureContainer(),
			txs:             smodels.NewTxsContainer(),
			balances:        smodels.NewAccountsContainer(),
		},
	}, nil
}

func (p *Parser) GetTaskExecutor(taskTitle string) (executor *smodels.Executor, err error) {
	switch taskTitle {
	case parserBaseTask:
		return &smodels.Executor{
			ExecHeight: p.ParseBase,
			Save:       p.Save,
		}, nil
	case parserBalancesSnapshotTask:
		return &smodels.Executor{
			ExecHeight: p.ParseBalancesSnapshot,
			Save:       p.Save,
		}, nil
	default:
		return nil, fmt.Errorf("executor %s not found", taskTitle)
	}
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
		p.container.blocks.Flush()
	}

	if !p.container.blockSignatures.IsEmpty() {
		tm := time.Now()
		err = p.dao.CreateBlockSignatures(p.container.blockSignatures.Signatures())
		if err != nil {
			return fmt.Errorf("dao.CreateBlockSignatures: %s", err.Error())
		}
		log.Print("Save time Signatures: ", time.Since(tm))

		p.container.blockSignatures.Flush()
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

		p.container.txs.Flush()
	}

	if !p.container.balances.IsEmpty() {
		tm := time.Now()
		err = p.dao.CreateAccountBalances(p.container.balances.Balances())
		if err != nil {
			return fmt.Errorf("dao.CreateAccountBalances: %s", err.Error())
		}

		log.Print("Save time Balances: ", time.Since(tm))

		p.container.balances.Flush()
	}

	return nil
}

func (p *Parser) saveAccounts() error {
	return nil
}

func (p *Parser) ParseBase(conn *grpcCommon.ClientConn, blockID uint64) error {

	parsTask, err := NewParserTask(p.ctx, conn, p.container)
	if err != nil {
		return err
	}

	err = parsTask.ParseBase(blockID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) ParseBalancesSnapshot(conn *grpcCommon.ClientConn, blockID uint64) error {

	parsTask, err := NewParserTask(p.ctx, conn, p.container)
	if err != nil {
		return err
	}

	err = parsTask.BalanceSnapshot(blockID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) ParseWatchBlock(block *consensusAPI.Block) error {
	parsTask, err := NewParserTask(p.ctx, p.conn, p.container)
	if err != nil {
		return err
	}

	err = parsTask.parseOasisBase(block, watcherFlag)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) getCustomHash(id string, seqNum uint64) (string, error) {
	key := append([]byte(id), []byte(strconv.Itoa(int(seqNum)))...)
	h, err := blake2b.New256(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

//Temp use hex
func hashHex(hash []byte) string {
	enc := make([]byte, len(hash)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], hash)
	return string(enc)
}
