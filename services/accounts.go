package services

import (
	"context"
	"github.com/oasislabs/oasis-core/go/common/crypto/signature"
	"github.com/oasislabs/oasis-core/go/staking/api"
	"log"
	"oasisTracker/smodels"
	"strings"
)

func (s *ServiceFacade) GetAccountInfo(accountID string) (sAcc smodels.Account, err error) {
	pb := signature.PublicKey{}

	err = pb.UnmarshalText([]byte(accountID))
	if err != nil {
		return sAcc, err
	}

	acc, err := s.nodeAPI.AccountInfo(context.Background(), &api.OwnerQuery{
		//Latest
		Height: 0,
		Owner:  pb,
	})
	if err != nil {
		log.Print(err)
	}

	log.Printf("%+v", acc)

	accType := "account"
	nodeAddress := ""

	for key, value := range acc.Escrow.StakeAccumulator.Claims {

		if len(value) == 1 {

			switch value[0] {
			case api.KindEntity:
				accType = value[0].String()
			case api.KindNodeValidator:
				splits := strings.Split(string(key), ".")
				if len(splits) == 3 {
					nodeAddress = splits[2]
				}

			}
		}
	}

	accountTime, err := s.dao.GetAccountTiming(accountID)
	if err != nil {
		return sAcc, err
	}

	liquidBalance := acc.General.Balance.ToBigInt().Uint64()
	escrowBalance := acc.Escrow.Active.Balance.ToBigInt().Uint64()

	return smodels.Account{
		Address:          accountID,
		LiquidBalance:    liquidBalance,
		EscrowBalance:    escrowBalance,
		DebondingBalance: acc.Escrow.Debonding.Balance.ToBigInt().Uint64(),
		TotalBalance:     liquidBalance + escrowBalance,
		CreatedAt:        accountTime.CreatedAt,
		LastActive:       accountTime.LastActive,
		Nonce:            acc.General.Nonce,
		Type:             accType,
		NodeAddress:      nodeAddress,
	}, nil
}
