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
	HighPrice  decimal.Decimal `json:"high_price"`
	LowPrice   decimal.Decimal `json:"low_price"`
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
		k.OpenPrice, k.HighPrice, k.LowPrice, k.ClosePrice,
	)
}

func (k *Kline) Update(kk *Kline) {
	k.Pair = kk.Pair
	k.Type = kk.Type
	k.Source = kk.Source
	k.OpenPrice = kk.OpenPrice
	k.ClosePrice = kk.ClosePrice
	k.HighPrice = kk.HighPrice
	k.LowPrice = kk.LowPrice
	k.Volume = kk.Volume
	k.Timestamp = kk.Timestamp
}
