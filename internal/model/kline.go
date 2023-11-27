package model

import "github.com/yanun0323/decimal"

type KlineSource int

const (
	KlineSourceUnknown KlineSource = iota
	KlineSourceBinance
	KlineSourceBitoPro
)

func (s KlineSource) IsUnknown() bool {
	return s == KlineSourceUnknown
}

type Kline struct {
	Pair   Pair
	Type   KlineType
	Source KlineSource

	OpenPrice  decimal.Decimal
	ClosePrice decimal.Decimal
	MaxPrice   decimal.Decimal
	MinPrice   decimal.Decimal
	Volume     decimal.Decimal

	Timestamp int64 /* end at (unix second) */
}
