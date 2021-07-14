package render

import (
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

func Rewards(rw []dmodels.Reward) []smodels.Reward {
	rewards := make([]smodels.Reward, len(rw))
	for i := range rw {
		rewards[i] = Reward(rw[i])
	}
	return rewards
}

func Reward(r dmodels.Reward) smodels.Reward {

	return smodels.Reward{
		Epoch:            r.Epoch,
		BlockLevel:       r.BlockLevel,
		ValidatorAddress: r.EntityAddress,
		Amount:           r.Amount,
		Type:             string(r.Type),
		CreatedAt:        r.CreatedAt.Unix(),
	}
}

func RewardStat(stat dmodels.RewardsStat) smodels.RewardStat {
	return smodels.RewardStat{
		AccountAddress: stat.AccountAddress,
		EntityAddress:  stat.EntityAddress,
		TotalAmount:    stat.TotalAmount,
		DayAmount:      stat.DayAmount,
		WeekAmount:     stat.WeekAmount,
		MonthAmount:    stat.MonthAmount,
	}
}
