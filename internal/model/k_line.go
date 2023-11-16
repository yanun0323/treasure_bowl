package model

import "github.com/yanun0323/decimal"

type KLineType uint8

const (
	K1m KLineType = iota
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

type KLine struct {
	Pair       Pair
	MaxPrice   decimal.Decimal
	MinPrice   decimal.Decimal
	OpenPrice  decimal.Decimal
	ClosePrice decimal.Decimal
	StartAt    uint64
	EndAt      uint64
	Timestamp  uint64
	Type       KLineType
}
