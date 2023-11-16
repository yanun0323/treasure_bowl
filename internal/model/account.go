package model

import (
	"main/internal/util"

	"github.com/yanun0323/decimal"
)

type Account struct {
	Balances util.SyncMap[string, Balance]
}

func NewAccount() *Account {
	return &Account{
		Balances: util.NewSyncMap[string, Balance](),
	}
}

type Balance struct {
	Available decimal.Decimal
	InTrade   decimal.Decimal
	Locked    decimal.Decimal
}
