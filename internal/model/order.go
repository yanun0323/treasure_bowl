package model

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

type Order struct {
	ID     string
	Pair   string
	Action OrderAction
	Type   OrderType

	LimitOrder
	MarketOrder
	StopLimitOrder
	TrailingStopOrder
	OCOOrder
}

type LimitOrder struct {
	Price  string
	Amount string
}

type MarketOrder struct {
	Amount string
	Total  string
}

type StopLimitOrder struct {
	Stop   string
	Limit  string
	Amount string
}

type TrailingStopOrder struct {
	TrailingDelta   string
	Limit           string
	Amount          string
	ActivationPrice string
}

type OCOOrder struct {
	Price  string
	Stop   string
	Limit  string
	Amount string
}
