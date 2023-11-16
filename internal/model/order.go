package model

import "github.com/yanun0323/decimal"

type OrderAction uint8

const (
	BUY OrderAction = iota
	SELL
)

type OrderType uint8

const (
	Limit OrderType = iota
	Market
	StopLimit
	TrailingStop
	OCO
)

type OrderStatus uint8

const (
	Pending OrderStatus = iota
	Created
	Canceled
	PartialComplete
	Complete
)

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
