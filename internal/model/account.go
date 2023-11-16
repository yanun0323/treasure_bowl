package model

import (
	"main/internal/util"

	"github.com/yanun0323/decimal"
)

type Balance struct {
	Available decimal.Decimal
	InTrade   decimal.Decimal
	Locked    decimal.Decimal
}

type Account struct {
	Balances util.SyncMap[string, Balance]
}

func NewAccount() Account {
	return Account{
		Balances: util.NewSyncMap[string, Balance](),
	}
}

func (acc *Account) CheckBalance(order Order) bool {
	if order.Action != BUY {
		return true
	}

	b, ok := acc.Balances.Load(order.Pair.Base())
	if !ok {
		return false
	}

	// TODO: Implement me
	_ = b
	return false
}
