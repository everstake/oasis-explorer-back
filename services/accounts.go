package services

import (
	"context"
	"oasisTracker/common/log"
	"oasisTracker/services/render"
	"oasisTracker/smodels"

	"go.uber.org/zap"

	"github.com/oasisprotocol/oasis-core/go/common/crypto/signature"
	"github.com/oasisprotocol/oasis-core/go/staking/api"
)

func (s *ServiceFacade) GetAccountInfo(accountID string) (sAcc smodels.Account, err error) {

	adr := api.NewAddress(signature.PublicKey{})
	err = adr.UnmarshalText([]byte(accountID))
	if err != nil {
		return sAcc, err
	}

	//Get last account state
	acc, err := s.nodeAPI.Account(context.Background(), &api.OwnerQuery{
		//Latest
		Height: 0,
		Owner:  adr,
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

			if *value[0].Global > kind {
				kind = *value[0].Global
			}
		}

		accType = kind.String()
	}

	liquidBalance := acc.General.Balance.ToBigInt().Uint64()
	escrowBalance := acc.Escrow.Active.Balance.ToBigInt().Uint64()

	//Get last account delegations state
	delegations, err := s.nodeAPI.DelegationInfosFor(context.Background(), &api.OwnerQuery{
		//Latest
		Height: 0,
		Owner:  adr,
	})
	if err != nil {
		return sAcc, err
	}

	var delegationsBalance, selfdelegation uint64
	for address, delegationInfo := range delegations {

		stakeBalance, err := delegationInfo.Pool.StakeForShares(&delegationInfo.Shares)
		if err != nil {
			log.Error("Somehow delegations rpc values is wrong", zap.Error(err))
			continue
		}

		//Self delegation
		if address.Equal(adr) {
			selfdelegation = stakeBalance.ToBigInt().Uint64()
		}
		delegationsBalance += stakeBalance.ToBigInt().Uint64()
	}

	//Get last account debonding delegations state
	debondingDelegations, err := s.nodeAPI.DebondingDelegationInfosFor(context.Background(), &api.OwnerQuery{
		//Latest
		Height: 0,
		Owner:  adr,
	})
	if err != nil {
		return sAcc, err
	}

	var debondingDelegationsBalance uint64
	for _, debDelegationList := range debondingDelegations {
		for _, value := range debDelegationList {

			debondingBalance, err := value.Pool.StakeForShares(&value.Shares)
			if err != nil {
				log.Error("Somehow debonding rpc values is wrong", zap.Error(err))
				continue
			}

			debondingDelegationsBalance += debondingBalance.ToBigInt().Uint64()
		}
	}

	sAcc = smodels.Account{
		Address:          accountID,
		LiquidBalance:    liquidBalance,
		EscrowBalance:    escrowBalance,
		DebondingBalance: acc.Escrow.Debonding.Balance.ToBigInt().Uint64(),

		DelegationsBalance:          delegationsBalance,
		DebondingDelegationsBalance: debondingDelegationsBalance,
		SelfDelegationBalance:       selfdelegation,
		TotalBalance:                liquidBalance + (escrowBalance - selfdelegation) + delegationsBalance + debondingDelegationsBalance,
		Type:                        accType,
		Nonce:                       &acc.General.Nonce,

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
		sAcc.EntityAddress = ent.EntityAddress

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

		sAcc.Type = smodels.AccountTypeValidator

		status := smodels.StatusActive
		if accType != api.KindNodeValidator.String() {
			status = smodels.StatusInActive
		}

		sAcc.Validator = &smodels.ValidatorInfo{
			CommissionScheduleRules: smodels.TestNetGenesis,
			Status:                  status,
			NodeAddress:             ent.Address,
			ConsensusAddress:        ent.ConsensusAddress,
			DepositorsCount:         depositorsCount,
			BlocksCount:             ent.BlocksCount,
			SignaturesCount:         ent.BlockSignaturesCount,
		}
	}

	return sAcc, nil
}

func (s *ServiceFacade) GetAccountList(listParams smodels.AccountListParams) (sAcc []smodels.AccountList, count uint64, err error) {

	count, err = s.dao.AccountsCount()
	if err != nil {
		return sAcc, 0, err
	}

	list, err := s.dao.GetAccountList(listParams)
	if err != nil {
		return sAcc, 0, err
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

	return render.AccountList(list), count, nil
}

func (s *ServiceFacade) GetAccountRewards(validatorID string, params smodels.CommonParams) ([]smodels.Reward, error) {
	rewards, err := s.dao.GetAccountRewards(validatorID, params)
	if err != nil {
		return nil, err
	}

	return render.Rewards(rewards), nil
}

func (s *ServiceFacade) GetAccountRewardsStat(validatorID string) (stat smodels.RewardStat, err error) {
	rewardsStat, err := s.dao.GetAccountRewardsStat(validatorID)
	if err != nil {
		return stat, err
	}

	return render.RewardStat(rewardsStat), nil
}
