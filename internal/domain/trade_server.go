package domain

import (
	"context"

	"main/internal/entity"
)

type TradeServer interface {
	Connect(context.Context) (entity.Account, <-chan entity.Order, error)
	Connected() bool
	Disconnect(context.Context) error
	IsSupported(entity.OrderType) bool
	PushOrder(context.Context, entity.Order) (entity.Account, error)
}
