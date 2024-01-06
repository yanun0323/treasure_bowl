package entity

import (
	"time"
)

type KlineType string

const (
	K1m  KlineType = "1m"
	K3m  KlineType = "3m"
	K5m  KlineType = "5m"
	K15m KlineType = "15m"
	K30m KlineType = "30m"
	K1h  KlineType = "1h"
	K2h  KlineType = "2h"
	K4h  KlineType = "4h"
	K6h  KlineType = "6h"
	K8h  KlineType = "8h"
	K12h KlineType = "12h"
	K1d  KlineType = "1d"
	K1w  KlineType = "1w"
	K1M  KlineType = "1M"
)

func (t KlineType) Validate() bool {
	switch t {
	case K1m, K3m, K5m, K15m, K30m, K1h, K2h, K4h, K6h, K8h, K12h, K1d, K1w, K1M:
		return true
	default:
		return false
	}
}

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

func (t KlineType) Spec() string {
	switch t {
	case K1m:
		return "5 * * * * *"
	case K3m:
		return "5 * * * * *"
	case K5m:
		return "5 * * * * *"
	case K15m:
		return "5 * * * * *"
	case K30m:
		return "5 * * * * *"
	case K1h:
		return "5 0 * * * *"
	case K2h:
		return "5 0 * * * *"
	case K4h:
		return "5 0 * * * *"
	case K6h:
		return "5 0 * * * *"
	case K8h:
		return "5 0 * * * *"
	case K12h:
		return "5 0 * * * *"
	case K1d:
		return "5 0 * * * *"
	case K1w:
		return "5 0 0 * * *"
	case K1M:
		return "5 0 0 * * *"

	}
	return "* * * * * *"
}

func (t KlineType) SpanDefault() Span {
	return t.Span(100)
}

func (t KlineType) Span(count int) Span {
	end := time.Now()
	return Span{
		Start: end.Add(-(t.Duration() * time.Duration(count))),
		End:   end,
	}
}

func (t KlineType) Duration() time.Duration {
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

// To returns the duration from start to end according to kline type
func (t KlineType) To(end int64) time.Duration {
	if t != K1M {
		return time.Duration(end) - t.Duration()
	}

	ed := time.Unix(end, 0)
	monthBegin := time.Date(ed.Year(), ed.Month(), 1, ed.Hour(), ed.Minute(), ed.Second(), ed.Nanosecond(), ed.Location())
	return ed.Sub(monthBegin.AddDate(0, 0, -1))
}
