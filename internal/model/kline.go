package model

import "github.com/yanun0323/decimal"

type Kline struct {
	Pair Pair
	Type KlineType

	OpenPrice  decimal.Decimal
	ClosePrice decimal.Decimal
	MaxPrice   decimal.Decimal
	MinPrice   decimal.Decimal
	Volume     decimal.Decimal

	Timestamp int64 /* end at (unix second) */
}
