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
	Pair   Pair        `json:"pair"`
	Type   KlineType   `json:"type"`
	Source KlineSource `json:"source"`

	OpenPrice  decimal.Decimal `json:"open_price"`
	ClosePrice decimal.Decimal `json:"close_price"`
	MaxPrice   decimal.Decimal `json:"max_price"`
	MinPrice   decimal.Decimal `json:"min_price"`
	Volume     decimal.Decimal `json:"volume"`

	Timestamp int64 `json:"timestamp"` /* end at (unix second) */
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
