package render

import (
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
	"strings"
)

func AccountList(accs []dmodels.AccountList) []smodels.AccountList {
	accounts := make([]smodels.AccountList, len(accs))
	for i := range accs {
		accounts[i] = AccountListElement(accs[i])
	}
	return accounts
}

func AccountListElement(a dmodels.AccountList) smodels.AccountList {

	return smodels.AccountList{
		Account:            a.Account,
		CreatedAt:          a.CreatedAt.Unix(),
		OperationsAmount:   a.OperationsAmount,
		GeneralBalance:     a.GeneralBalance,
		EscrowBalance:      a.EscrowBalanceActive,
		EscrowBalanceShare: a.EscrowBalanceShare,
		Delegate:           strings.Trim(a.Delegate, "\u0000"),
		Type:               a.Type,
	}
}
