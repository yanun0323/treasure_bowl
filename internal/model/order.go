package model

import "github.com/yanun0323/decimal"

const (
	_unknown = "Unknown"
)

// OrderAction
type OrderAction uint8

const (
	BUY OrderAction = iota
	SELL
	CANCEL
)

func (a OrderAction) String() string {
	switch a {
	case BUY:
		return "BUY"
	case SELL:
		return "SELL"
	case CANCEL:
		return "CANCEL"
	default:
		return _unknown
	}
}

// OrderType
type OrderType uint8

const (
	Limit OrderType = iota
	Market
	StopLimit
	TrailingStop
	OCO
)

func (t OrderType) String() string {
	switch t {
	case Limit:
		return "Limit"
	case Market:
		return "Market"
	case StopLimit:
		return "StopLimit"
	case TrailingStop:
		return "TrailingStop"
	case OCO:
		return "OCO"
	default:
		return _unknown
	}
}

// OrderStatus
type OrderStatus uint8

const (
	Pending OrderStatus = iota
	Created
	Canceled
	PartialComplete
	Complete
)

func (s OrderStatus) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Created:
		return "Created"
	case Canceled:
		return "Canceled"
	case PartialComplete:
		return "PartialComplete"
	case Complete:
		return "Complete"
	default:
		return _unknown
	}
}

// Order
type Order struct {
	ID     string
	Pair   Pair
	Action OrderAction
	Type   OrderType
	Status OrderStatus

	LimitOrder
	MarketOrder
	StopLimitOrder
	TrailingStopOrder
	OCOOrder
}

type LimitOrder struct {
	Price  decimal.Decimal
	Amount decimal.Decimal
}

type MarketOrder struct {
	Amount decimal.Decimal
	Total  decimal.Decimal
}

type StopLimitOrder struct {
	Stop   decimal.Decimal
	Limit  decimal.Decimal
	Amount decimal.Decimal
}

type TrailingStopOrder struct {
	TrailingDelta   decimal.Decimal
	Limit           decimal.Decimal
	Amount          decimal.Decimal
	ActivationPrice decimal.Decimal
}

type OCOOrder struct {
	Price  decimal.Decimal
	Stop   decimal.Decimal
	Limit  decimal.Decimal
	Amount decimal.Decimal
}
