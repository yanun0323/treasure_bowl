package entity

import (
	"fmt"
	"reflect"

	"github.com/yanun0323/decimal"
)

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

func (k *Kline) IsTypeEqual(kk *Kline) bool {
	return k.Pair == kk.Pair && k.Type == kk.Type && k.Source == kk.Source && k.Timestamp == kk.Timestamp
}

func (k *Kline) IsEqual(kk *Kline) bool {
	if !k.IsTypeEqual(kk) {
		return false
	}
	return reflect.DeepEqual(*k, *kk)
}

func (k Kline) String() string {
	return fmt.Sprintf("%s %s %d, ohlc: %s %s %s %s ",
		k.Pair.Uppercase("_"), k.Type.String(), k.Timestamp,
		k.OpenPrice, k.MaxPrice, k.MinPrice, k.ClosePrice,
	)
}

func (k *Kline) Update(kk *Kline) {
	k.Pair = kk.Pair
	k.Type = kk.Type
	k.Source = kk.Source
	k.OpenPrice = kk.OpenPrice
	k.ClosePrice = kk.ClosePrice
	k.MaxPrice = kk.MaxPrice
	k.MinPrice = kk.MinPrice
	k.Volume = kk.Volume
	k.Timestamp = kk.Timestamp
}
