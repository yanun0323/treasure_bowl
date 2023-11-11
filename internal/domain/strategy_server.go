package domain

import (
	"context"
	"main/internal/model"
)

type StrategyServer interface {
	Connect(ctx context.Context) (<-chan model.Order, error)
	Disconnect(ctx context.Context) error

	PushPrices(ctx context.Context, prices ...model.Price)
	PushAssets(ctx context.Context, accounts ...model.Account)
	PushOrders(ctx context.Context, orders ...model.Order)
	PushSupportedOrderTypes(ctx context.Context, types ...model.OrderType)
}
