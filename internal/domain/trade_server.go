package domain

import (
	"context"

	"main/internal/model"
)

type TradeServer interface {
	Connect(context.Context) (model.Account, <-chan model.Order, error)
	Connected() bool
	Disconnect(context.Context) error
	IsSupported(model.OrderType) bool
	PushOrder(context.Context, model.Order) (model.Account, error)
}
