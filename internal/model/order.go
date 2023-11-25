package model

import (
	"github.com/yanun0323/decimal"
)

const (
	_unknown = "Unknown"
)

// OrderAction requirement for pushing order to order server
type OrderAction uint8

const (
	None OrderAction = iota /* 'None' set when this order comes from the order server */
	BUY
	SELL
	BUY_CANCEL
	SELL_CANCEL
)

func (a OrderAction) String() string {
	switch a {
	case None:
		return "NONE"
	case BUY:
		return "BUY"
	case SELL:
		return "SELL"
	case BUY_CANCEL:
		return "BUY_CANCEL"
	case SELL_CANCEL:
		return "SELL_CANCEL"
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
	Complete
	PartialComplete
	Canceled
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
	ID        string
	Pair      Pair
	Action    OrderAction
	Type      OrderType
	Status    OrderStatus
	Timestamp int64 /* unix second */
	Price     decimal.Decimal
	Amount    Amount

	StopLimitOrder
	TrailingStopOrder
	OCOOrder
}

func (order *Order) GetTotal() decimal.Decimal {
	return order.Amount.Total
}

type Amount struct {
	Total  decimal.Decimal
	Deal   decimal.Decimal
	Remain decimal.Decimal
}

type StopLimitOrder struct {
	Stop decimal.Decimal
}

type TrailingStopOrder struct {
	TrailingDelta   decimal.Decimal
	Limit           decimal.Decimal
	ActivationPrice decimal.Decimal
}

type OCOOrder struct {
	Stop  decimal.Decimal
	Limit decimal.Decimal
}
