package services

import (
	"context"
	"github.com/oasislabs/oasis-core/go/common/crypto/signature"
	"github.com/oasislabs/oasis-core/go/staking/api"
	"oasisTracker/services/render"
	"oasisTracker/smodels"
)

func (s *ServiceFacade) GetAccountInfo(accountID string) (sAcc smodels.Account, err error) {
	pb := signature.PublicKey{}

	err = pb.UnmarshalText([]byte(accountID))
	if err != nil {
		return sAcc, err
	}

	//Get last account state
	acc, err := s.nodeAPI.AccountInfo(context.Background(), &api.OwnerQuery{
		//Latest
		Height: 0,
		Owner:  pb,
	})
	if err != nil {
		return sAcc, err
	}

	//Get account Create LastActive time based on txs
	accountTime, err := s.dao.GetAccountTiming(accountID)
	if err != nil {
		return sAcc, err
	}

	accType := smodels.AccountTypeAccount

	//Account can be entity but doesn't have validator node
	if len(acc.Escrow.StakeAccumulator.Claims) > 0 {
		var kind api.ThresholdKind
		for _, value := range acc.Escrow.StakeAccumulator.Claims {
			if len(value) < 1 {
				continue
			}

			if value[0] > kind {
				kind = value[0]
			}
		}

		accType = kind.String()
	}

	liquidBalance := acc.General.Balance.ToBigInt().Uint64()
	escrowBalance := acc.Escrow.Active.Balance.ToBigInt().Uint64()

	sAcc = smodels.Account{
		Address:          accountID,
		LiquidBalance:    liquidBalance,
		EscrowBalance:    escrowBalance,
		DebondingBalance: acc.Escrow.Debonding.Balance.ToBigInt().Uint64(),
		TotalBalance:     liquidBalance + escrowBalance,
		Type:             accType,
		Nonce:            &acc.General.Nonce,

		CreatedAt:  accountTime.CreatedAt,
		LastActive: accountTime.LastActive,
	}

	//Check all account because node addresses are displayed only on Entity address
	resp, err := s.dao.GetAccountValidatorInfo(accountID)
	if err != nil {
		return sAcc, err
	}

	switch {
	//Node account
	case resp.IsNode(accountID):
		sAcc.EntityAddress = resp.GetEntityAddress()
		sAcc.Type = smodels.AccountTypeNode
	//Entity account
	case resp.IsEntity(accountID):
		ent := resp.GetEntity()
		sAcc.EntityAddress = ent.EntityID

		depositorsCount, err := s.dao.GetEntityActiveDepositorsCount(accountID)
		if err != nil {
			return sAcc, err
		}

		if ent.CreatedTime.Before(sAcc.CreatedAt) {
			sAcc.CreatedAt = ent.CreatedTime
		}

		lastActive := ent.GetLastActiveTime()
		if lastActive.After(sAcc.LastActive) {
			sAcc.LastActive = lastActive
		}

		sAcc.Type = "validator"

		status := "active"
		if accType != api.KindNodeValidator.String() {
			status = "inactive"
		}

		sAcc.Validator = &smodels.Validator{
			CommissionScheduleRules: smodels.TestNetGenesis,
			Status:                  status,
			NodeAddress:             ent.NodeID,
			ConsensusAddress:        ent.ConsensusAddress,
			DepositorsCount:         depositorsCount,
			BlocksCount:             ent.BlocksCount,
			SignaturesCount:         ent.BlockSignaturesCount,
		}
	}

	return sAcc, nil
}

func (s *ServiceFacade) GetAccountList(listParams smodels.AccountListParams) (sAcc []smodels.AccountList, err error) {

	list, err := s.dao.GetAccountList(listParams)
	if err != nil {
		return sAcc, err
	}

	for i := range list {
		accountType := smodels.AccountTypeAccount
		switch {
		case list[i].NodeRegisterBlock > 0:
			accountType = smodels.AccountTypeNode
		case list[i].EntityRegisterBlock > 0:
			accountType = smodels.AccountTypeEntity
		}

		list[i].Type = accountType
	}

	return render.AccountList(list), nil
}
