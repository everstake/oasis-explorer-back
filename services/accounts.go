package services

import (
	"context"
	"github.com/oasislabs/oasis-core/go/common/crypto/signature"
	"github.com/oasislabs/oasis-core/go/staking/api"
	"github.com/wedancedalot/decimal"
	"log"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
	"strings"
)

func (s *ServiceFacade) GetAccountInfo(accountID string) (smodels.Account, error) {
	pb := signature.PublicKey{}

	err := pb.UnmarshalText([]byte(accountID))
	if err != nil {

	}

	acc, err := s.nodeAPI.AccountInfo(context.Background(), &api.OwnerQuery{
		Height: 201978,
		Owner:  pb,
	})
	if err != nil {
		log.Print(err)
	}

	log.Printf("%+v", acc)

	log.Print(len(acc.Escrow.StakeAccumulator.Claims))

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
		log.Print(err)
	}

	liquidBalance := decimal.NewFromBigInt(acc.General.Balance.ToBigInt(), int32(-dmodels.Precision))
	escrowBalance := decimal.NewFromBigInt(acc.Escrow.Active.Balance.ToBigInt(), int32(-dmodels.Precision))

	return smodels.Account{
		Address:          accountID,
		LiquidBalance:    liquidBalance,
		EscrowBalance:    escrowBalance,
		DebondingBalance: decimal.NewFromBigInt(acc.Escrow.Debonding.Balance.ToBigInt(), int32(-dmodels.Precision)),
		TotalBalance:     decimal.Sum(liquidBalance, escrowBalance),
		CreatedAt:        accountTime.CreatedAt,
		LastActive:       accountTime.LastActive,
		Nonce:            acc.General.Nonce,
		Type:             accType,
		NodeAddress:      nodeAddress,
	}, nil
}
