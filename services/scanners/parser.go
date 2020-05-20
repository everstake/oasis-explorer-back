package scanners

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/fxamacker/cbor/v2"
	"github.com/oasislabs/oasis-core/go/common/grpc"
	consensusAPI "github.com/oasislabs/oasis-core/go/consensus/api"
	"github.com/tendermint/tendermint/types"
	"golang.org/x/crypto/blake2b"
	grpcCommon "google.golang.org/grpc"
	"log"
	"oasisTracker/conf"
	"oasisTracker/dmodels"
	"oasisTracker/dmodels/oasis"
	"oasisTracker/smodels"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

const (
	precision          = 6
	saveBatch          = 200
	saveAddressesBatch = 50

	parserBaseTask        = "base"
	parserSignaturesTask  = "signatures"
	parseTransactionsTask = "transactions"

	parseBlock ParseFlag = iota
	parseSig
	parseFullBlock
)

type (
	ParseFlag uint32
	Parser    struct {
		ctx context.Context

		dao DAO
		api consensusAPI.ClientBackend

		blocks          *smodels.BlocksContainer
		blockSignatures *smodels.BlockSignatureContainer
		txs             *smodels.TxsContainer
	}
	DAO interface {
		CreateBlocks(blocks []dmodels.Block) error
		CreateBlockSignatures(sig []dmodels.BlockSignature) error
		//CreateAccounts(accounts []interface{}) error
		CreateTransfers(transfers []dmodels.Transaction) error
	}
)

func NewParser(ctx context.Context, cfg conf.Scanner, tezosdDAO interface{}) (*Parser, error) {
	grpcConn, err := grpc.Dial(cfg.NodeConfig, grpcCommon.WithInsecure())
	if err != nil {
		log.Print(err)
	}

	cAPI := consensusAPI.NewConsensusClient(grpcConn)

	d, ok := tezosdDAO.(DAO)
	if !ok {
		return nil, fmt.Errorf("can`t cast to oasis DAO")
	}
	return &Parser{
		ctx:             ctx,
		api:             cAPI,
		dao:             d,
		blocks:          smodels.NewBlocksContainer(),
		blockSignatures: smodels.NewBlockSignatureContainer(),
		txs:             smodels.NewTxsContainer(),
	}, nil
}

func (p *Parser) GetTaskExecutor(taskTitle string) (executor *smodels.Executor, err error) {
	switch taskTitle {
	case parserBaseTask:
		return &smodels.Executor{
			ExecHeight: p.ParseBlockData,
			Save:       p.Save,
		}, nil
	case parserSignaturesTask:
		return &smodels.Executor{
			ExecHeight: p.ParseBlockSignatures,
			Save:       p.Save,
		}, nil
	case parseTransactionsTask:
		return &smodels.Executor{
			ExecHeight: p.ParseBlockTransactions,
			Save:       p.Save,
		}, nil
	default:
		return nil, fmt.Errorf("executor %s not found", taskTitle)
	}
}

func (p *Parser) Save() (err error) {

	if !p.blocks.IsEmpty() {
		err := p.dao.CreateBlocks(p.blocks.Blocks())
		if err != nil {
			return fmt.Errorf("dao.CreateBlocks: %s", err.Error())
		}

		p.blocks.Flush()
	}

	if !p.blockSignatures.IsEmpty() {
		tm := time.Now()
		err = p.dao.CreateBlockSignatures(p.blockSignatures.Signatures())
		if err != nil {
			return fmt.Errorf("dao.CreateBlockSignatures: %s", err.Error())
		}
		log.Print("Save time Signatures: ", time.Since(tm))

		p.blockSignatures.Flush()
	}

	if !p.txs.IsEmpty() {
		err = p.dao.CreateTransfers(p.txs.Txs())
		if err != nil {
			return fmt.Errorf("dao.CreateTransfers: %s", err.Error())
		}

		p.txs.Flush()
	}

	return nil
}

func (p *Parser) saveAccounts() error {
	return nil
}

func (p *Parser) ParseFullBlock(blockID uint64) (err error) {
	err = p.ParseBase(blockID, parseFullBlock)
	if err != nil {
		return err
	}

	err = p.ParseBlockTransactions(blockID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) ParseBlockSignatures(blockID uint64) (err error) {
	err = p.ParseBase(blockID, parseSig)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) ParseBlockData(blockID uint64) (err error) {
	err = p.ParseBase(blockID, parseBlock)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) ParseBlockTransactions(blockID uint64) (err error) {
	blockData, err := p.api.GetBlock(p.ctx, int64(blockID))
	if err != nil {
		return fmt.Errorf("api.Block.Get: %s", err.Error())
	}

	txs, err := p.api.GetTransactions(p.ctx, int64(blockID))
	if err != nil {
		return err
	}

	for key := range txs {
		tx := oasis.TxRaw{}

		err = cbor.Unmarshal(txs[key], &tx)
		if err != nil {
			return err
		}

		raw := oasis.UntrustedRawValue{}
		err = cbor.Unmarshal(tx.UntrustedRawValue, &raw)
		if err != nil {
			return err
		}

		p.txs.Add([]dmodels.Transaction{{
			BlockLevel:    blockID,
			Hash:          hex.EncodeToString(types.Tx(txs[key]).Hash()),
			Time:          blockData.Time,
			Amount:        raw.Body.Transfer.Tokens.ToBigInt().Uint64(),
			EscrowAmount:  raw.Body.EscrowTx.Tokens.ToBigInt().Uint64(),
			EscrowAccount: raw.Body.EscrowTx.Account.String(),
			Type:          dmodels.TransactionType(raw.Method),
			Sender:        tx.Signature.PublicKey.String(),
			Receiver:      raw.Body.To.String(),
			Nonce:         raw.Nonce,
			Fee:           raw.Fee.Amount.ToBigInt().Uint64(),
			GasLimit:      uint64(raw.Fee.Gas),
			GasPrice:      raw.Fee.GasPrice().ToBigInt().Uint64(),
		}})
	}

	return nil
}

func (p *Parser) ParseBase(blockID uint64, flag ParseFlag) error {
	blockData, err := p.api.GetBlock(p.ctx, int64(blockID))
	if err != nil {
		return fmt.Errorf("api.Block.Get: %s", err.Error())
	}

	b := oasis.Block{}
	err = cbor.Unmarshal(blockData.Meta, &b)
	if err != nil {
		return err
	}

	b.Hash = blockData.Hash

	var pipes []func(data oasis.Block) error

	switch flag {
	case parseFullBlock:
		fallthrough
	case parseBlock:
		pipes = append(pipes, p.parseBlock)
	case parseSig:
		pipes = append(pipes, p.parseBlockSignatures)
	default:
		return fmt.Errorf("Unknown flag: %d", flag)
	}

	for _, pipe := range pipes {
		err = pipe(b)
		if err != nil {
			funcName := runtime.FuncForPC(reflect.ValueOf(pipe).Pointer()).Name()
			return fmt.Errorf("%s (block:%d): %s", funcName, blockData.Height, err.Error())
		}
	}
	return nil
}

func (p *Parser) parseBlock(block oasis.Block) error {
	epoch, err := p.api.GetEpoch(p.ctx, block.Header.Height)
	if err != nil {
		return err
	}

	p.blocks.Add([]dmodels.Block{{
		Height:          uint64(block.Header.Height),
		Hash:            hex.EncodeToString(block.Hash),
		CreatedAt:       block.Header.Time,
		Epoch:           uint64(epoch),
		ProposerAddress: block.Header.ProposerAddress.String(),
		ValidatorHash:   block.Header.ValidatorsHash.String(),
	}})

	return nil
}

func (p *Parser) parseBlockSignatures(block oasis.Block) error {

	blockSignatures := make([]dmodels.BlockSignature, 0, len(block.LastCommit.Signatures))

	for key := range block.LastCommit.Signatures {
		timestamp := block.LastCommit.Signatures[key].Timestamp
		if timestamp.IsZero() {
			timestamp = block.Header.Time
		}

		//Use block timestamp if signature time is zero
		blockSignatures = append(blockSignatures, dmodels.BlockSignature{
			BlockHeight:      block.Header.Height,
			Timestamp:        timestamp,
			BlockIDFlag:      block.LastCommit.Signatures[key].BlockIDFlag,
			ValidatorAddress: block.LastCommit.Signatures[key].ValidatorAddress.String(),
			Signature:        hex.EncodeToString(block.LastCommit.Signatures[key].Signature),
		})

	}

	p.blockSignatures.Add(blockSignatures)

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
