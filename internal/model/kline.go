package model

import "github.com/yanun0323/decimal"

type Kline struct {
	Pair       Pair
	MaxPrice   decimal.Decimal
	MinPrice   decimal.Decimal
	OpenPrice  decimal.Decimal
	ClosePrice decimal.Decimal
	Volume     decimal.Decimal
	Type       KlineType
	Timestamp  int64 /* end at (unix second) */
}
