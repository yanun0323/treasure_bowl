package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/yanun0323/decimal"
)

type Account interface {
	Available(currency string) decimal.Decimal
	Unavailable(currency string) decimal.Decimal
	Timestamp() int64 /* unix second */
}

type Balance struct {
	Available   decimal.Decimal
	Unavailable decimal.Decimal
}

type account struct {
	balance   map[string]Balance
	timestamp int64
}

func NewAccount(balances ...map[string]Balance) Account {
	m := map[string]Balance{}
	if len(balances) != 0 {
		m = balances[0]
	}

	return &account{
		balance:   m,
		timestamp: time.Now().Unix(),
	}
}

func (a account) Available(currency string) decimal.Decimal {
	return a.balance[strings.ToUpper(currency)].Available
}

func (a account) Unavailable(currency string) decimal.Decimal {
	return a.balance[strings.ToUpper(currency)].Unavailable
}

func (a account) Timestamp() int64 {
	return a.timestamp
}

func (a *account) Trade(currency string, amount decimal.Decimal) error {
	c := strings.ToUpper(currency)
	b := a.balance[c]
	if !b.Available.GreaterOrEqual(amount) {
		return errors.New(fmt.Sprintf("%s insufficient balance, need %s but left %s ",
			c,
			amount,
			b.Available,
		))
	}

	b.Available = b.Available.Sub(amount)
	b.Unavailable = b.Unavailable.Add(amount)
	a.balance[c] = b
	return nil
}
