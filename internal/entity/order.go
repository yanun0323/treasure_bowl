package entity

import (
	"github.com/pkg/errors"
	"github.com/yanun0323/decimal"
)

const (
	_unknown = "UNKNOWN"
)

// OrderAction requirement for pushing order to order server
type OrderAction uint8

const (
	OrderActionNone OrderAction = iota /* 'None' set when this order comes from the order server */
	OrderActionBuy
	OrderActionSell
	OrderActionCancelBuy
	OrderActionCancelSell
)

func (a OrderAction) IsNone() bool {
	return a == OrderActionNone
}

func (a OrderAction) String() string {
	switch a {
	case OrderActionNone:
		return "NONE"
	case OrderActionBuy:
		return "BUY"
	case OrderActionSell:
		return "SELL"
	case OrderActionCancelBuy:
		return "CANCEL_BUY"
	case OrderActionCancelSell:
		return "CANCEL_SELL"
	default:
		return _unknown
	}
}

// OrderType
type OrderType uint8

const (
	OrderTypeUnknown OrderType = iota
	OrderTypeLimit
	OrderTypeMarket
	OrderTypeStopLimit
	OrderTypeTrailingStop
	OrderTypeOCO
)

func (t OrderType) IsUnknown() bool {
	return t == OrderTypeUnknown
}

func (t OrderType) String() string {
	switch t {
	case OrderTypeLimit:
		return "Limit"
	case OrderTypeMarket:
		return "Market"
	case OrderTypeStopLimit:
		return "StopLimit"
	case OrderTypeTrailingStop:
		return "TrailingStop"
	case OrderTypeOCO:
		return "OCO"
	default:
		return _unknown
	}
}

// OrderStatus
type OrderStatus uint8

const (
	OrderStatusUnknown OrderStatus = iota
	OrderStatusPending
	OrderStatusCreated
	OrderStatusComplete
	OrderStatusPartialComplete
	OrderStatusCanceled
)

func (s OrderStatus) IsUnknown() bool {
	return s == OrderStatusUnknown
}

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "Pending"
	case OrderStatusCreated:
		return "Created"
	case OrderStatusCanceled:
		return "Canceled"
	case OrderStatusPartialComplete:
		return "PartialComplete"
	case OrderStatusComplete:
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

func (o *Order) ValidatePushingOrder() error {
	if o.Action.IsNone() {
		return errors.New("invalid order action: None")
	}

	if o.Type.IsUnknown() {
		return errors.New("unknown order type")
	}

	return nil
}

func (o *Order) GetTotal() decimal.Decimal {
	return o.Amount.Total
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
