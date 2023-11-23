package domain

import (
	"context"

	"main/internal/model"
)

type StrategyServer interface {
	Connect(ctx context.Context) (<-chan model.Order, error)
	Disconnect(ctx context.Context) error

	PushKline(ctx context.Context, prices ...model.Kline)
	PushAsset(ctx context.Context, accounts ...model.Account)
	PushOrder(ctx context.Context, orders ...model.Order)
	PushSupportedOrderType(ctx context.Context, types ...model.OrderType)
}
