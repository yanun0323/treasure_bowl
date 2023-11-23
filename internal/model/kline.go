package model

import "github.com/yanun0323/decimal"

type KlineType uint8

const (
	K1m KlineType = iota
	K3m
	K5m
	K15m
	K30m
	K1h
	K2h
	K4h
	K6h
	K8h
	K12h
	K1d
	K3d
	K1w
	K1M
)

type Kline struct {
	Pair       Pair
	MaxPrice   decimal.Decimal
	MinPrice   decimal.Decimal
	OpenPrice  decimal.Decimal
	ClosePrice decimal.Decimal
	Type       KlineType
	Timestamp  int64
}
