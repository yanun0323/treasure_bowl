package model

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/yanun0323/decimal"

	"main/internal/util"
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

func (acc *Account) MoveToInTrade(currency string, amount decimal.Decimal) error {
	b, ok := acc.Balances.Load(currency)
	if !ok {
		return errors.New(fmt.Sprintf("currency %s has no balance", currency))
	}

	if b.Available.Less(amount) {
		return errors.New(fmt.Sprintf("available balance of currency %s is not enough. need: %s, actual: %s", currency, amount, b.Available))
	}

	b.Available = b.Available.Sub(amount)
	b.InTrade = b.InTrade.Add(amount)
	acc.Balances.Store(currency, b)

	return nil
}

func (acc *Account) MoveToAvailable(currency string, amount decimal.Decimal) error {
	b, ok := acc.Balances.Load(currency)
	if !ok {
		return errors.New(fmt.Sprintf("currency %s has no balance", currency))
	}

	if b.InTrade.Less(amount) {
		return errors.New(fmt.Sprintf("trading balance of currency %s is not enough. need: %s, actual: %s", currency, amount, b.InTrade))
	}

	b.InTrade = b.InTrade.Sub(amount)
	b.Available = b.Available.Add(amount)
	acc.Balances.Store(currency, b)

	return nil
}
