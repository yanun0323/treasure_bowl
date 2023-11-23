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
	K1w
	K1M
)

func (t KlineType) CronSpec() string {
	switch t {
	case K1m:
		return "0 * * * * *"
	case K3m:
		return "0 */3 * * * *"
	case K5m:
		return "0 */5 * * * *"
	case K15m:
		return "0 */15 * * * *"
	case K30m:
		return "0 */30 * * * *"
	case K1h:
		return "0 0 * * * *"
	case K2h:
		return "0 0 */2 * * *"
	case K4h:
		return "0 0 */4 * * *"
	case K6h:
		return "0 0 */6 * * *"
	case K8h:
		return "0 0 */8 * * *"
	case K12h:
		return "0 0 */12 * * *"
	case K1d:
		return "0 0 0 * * *"
	case K1w:
		return "0 0 0 0 0 *"
	case K1M:
		return "0 0 0 0 * *"

	}
	return "0 * * * * *"
}

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
