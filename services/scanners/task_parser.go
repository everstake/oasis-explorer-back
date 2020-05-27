package scanners

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/fxamacker/cbor/v2"
	consensusAPI "github.com/oasislabs/oasis-core/go/consensus/api"
	"github.com/tendermint/tendermint/types"
	"github.com/wedancedalot/decimal"
	"google.golang.org/grpc"
	"oasisTracker/dmodels"
	"oasisTracker/dmodels/oasis"
	"reflect"
	"runtime"
)

//Struct for direct work with connection from worker
type ParserTask struct {
	ctx             context.Context
	api             consensusAPI.ClientBackend
	parserContainer *ParseContainer
}

func NewParserTask(ctx context.Context, conn *grpc.ClientConn, parserContainer *ParseContainer) (*ParserTask, error) {
	api := consensusAPI.NewConsensusClient(conn)

	return &ParserTask{ctx: ctx, api: api, parserContainer: parserContainer}, nil
}

func (p *ParserTask) ParseBase(blockID uint64) error {
	blockData, err := p.api.GetBlock(p.ctx, int64(blockID))
	if err != nil {
		return fmt.Errorf("api.Block.Get: %s", err.Error())
	}

	err = p.ParseOasisBase(blockData)
	if err != nil {
		return fmt.Errorf("p.parseOasisBase: %s", err.Error())
	}

	return nil
}

func (p *ParserTask) ParseOasisBase(blockData *consensusAPI.Block) (err error) {

	b := oasis.Block{}
	err = cbor.Unmarshal(blockData.Meta, &b)
	if err != nil {
		return err
	}

	b.Hash = blockData.Hash

	pipes := []func(data oasis.Block) error{
		p.parseBlock,
		p.parseBlockSignatures,
		p.parseBlockTransactions,
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

func (p *ParserTask) parseBlock(block oasis.Block) error {
	epoch, err := p.api.GetEpoch(p.ctx, block.Header.Height)
	if err != nil {
		return err
	}

	p.parserContainer.blocks.Add([]dmodels.Block{{
		Height:          uint64(block.Header.Height),
		Hash:            hex.EncodeToString(block.Hash),
		CreatedAt:       block.Header.Time,
		Epoch:           uint64(epoch),
		ProposerAddress: block.Header.ProposerAddress.String(),
		ValidatorHash:   block.Header.ValidatorsHash.String(),
	}})

	return nil
}

func (p *ParserTask) parseBlockSignatures(block oasis.Block) error {

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

	p.parserContainer.blockSignatures.Add(blockSignatures)

	return nil
}

func (p *ParserTask) parseBlockTransactions(block oasis.Block) (err error) {
	txs, err := p.api.GetTransactions(p.ctx, block.Header.Height)
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

		txType, err := dmodels.NewTransactionType(raw.Method)
		if err != nil {
			return err
		}

		p.parserContainer.txs.Add([]dmodels.Transaction{{
			BlockLevel:          uint64(block.Header.Height),
			BlockHash:           block.Hash.String(),
			Hash:                hex.EncodeToString(types.Tx(txs[key]).Hash()),
			Time:                block.Header.Time,
			Amount:              decimal.NewFromBigInt(raw.Body.Transfer.Tokens.ToBigInt(), -int32(dmodels.Precision)).String(),
			EscrowAmount:        decimal.NewFromBigInt(raw.Body.EscrowTx.Tokens.ToBigInt(), -int32(dmodels.Precision)).String(),
			EscrowReclaimAmount: decimal.NewFromBigInt(raw.Body.EscrowTx.Shares.ToBigInt(), -int32(dmodels.Precision)).String(),
			EscrowAccount:       raw.Body.EscrowTx.Account.String(),
			Type:                txType.Type(),
			Sender:              tx.Signature.PublicKey.String(),
			Receiver:            raw.Body.To.String(),
			Nonce:               raw.Nonce,
			Fee:                 raw.Fee.Amount.ToBigInt().Uint64(),
			GasLimit:            uint64(raw.Fee.Gas),
			GasPrice:            raw.Fee.GasPrice().ToBigInt().Uint64(),
		}})
	}

	return nil
}
