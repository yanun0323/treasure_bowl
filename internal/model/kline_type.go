package model

import (
	"time"
)

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

func (t KlineType) String() string {
	switch t {
	case K1m:
		return "1m"
	case K3m:
		return "3m"
	case K5m:
		return "5m"
	case K15m:
		return "15m"
	case K30m:
		return "30m"
	case K1h:
		return "1h"
	case K2h:
		return "2h"
	case K4h:
		return "4h"
	case K6h:
		return "6h"
	case K8h:
		return "8h"
	case K12h:
		return "12h"
	case K1d:
		return "1d"
	case K1w:
		return "1w"
	case K1M:
		return "1M"

	}
	return ""
}

func (t KlineType) duration() time.Duration {
	switch t {
	case K1m:
		return time.Minute
	case K3m:
		return 3 * time.Minute
	case K5m:
		return 5 * time.Minute
	case K15m:
		return 15 * time.Minute
	case K30m:
		return 30 * time.Minute
	case K1h:
		return time.Hour
	case K2h:
		return 2 * time.Hour
	case K4h:
		return 4 * time.Hour
	case K6h:
		return 6 * time.Hour
	case K8h:
		return 8 * time.Hour
	case K12h:
		return 12 * time.Hour
	case K1d:
		return 24 * time.Hour
	case K1w:
		return 24 * 7 * time.Hour
	case K1M:
		return 24 * 7 * 30 * time.Hour
	}
	return time.Minute
}

// Duration return the duration from start to end according to kline type
func (t KlineType) Duration(end int64) time.Duration {
	if t != K1M {
		return time.Duration(end) - t.duration()
	}

	ed := time.Unix(end, 0)
	monthBegin := time.Date(ed.Year(), ed.Month(), 1, ed.Hour(), ed.Minute(), ed.Second(), ed.Nanosecond(), ed.Location())
	return ed.Sub(monthBegin.AddDate(0, 0, -1))
}
