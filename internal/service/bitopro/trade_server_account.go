package bitopro

import (
	"strings"
	"time"

	"main/internal/model"

	"github.com/bitoex/bitopro-api-go/pkg/bitopro"
	"github.com/pkg/errors"
	"github.com/yanun0323/decimal"
)

var (
	/* Check Interface Implement */
	_ model.Account = (*account)(nil)
)

type account struct {
	Data *bitopro.Account
	TS   int64 /* unix second */
}

// convAccount converts bitopro account into model account.
func convAccount(a *bitopro.Account) (model.Account, error) {
	if a == nil {
		return nil, errors.New("nil account data")
	}

	if len(a.Error) != 0 {
		return nil, errors.New("converts account: " + a.Error)
	}
	return &account{
		Data: a,
		TS:   time.Now().Unix(),
	}, nil
}

func (a account) Available(currency string) decimal.Decimal {
	for _, d := range a.Data.Data {
		if strings.EqualFold(currency, d.Currency) {
			return decimal.Require(d.Available)
		}
	}
	return decimal.Zero()
}

func (a account) Unavailable(currency string) decimal.Decimal {
	for _, d := range a.Data.Data {
		if strings.EqualFold(currency, d.Currency) {
			return decimal.Require(d.Amount).Sub(decimal.Require(d.Available))
		}
	}
	return decimal.Zero()
}

func (a account) Timestamp() int64 /* unix second */ {
	return a.TS
}
