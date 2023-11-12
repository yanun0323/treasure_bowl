package domain

import (
	"context"
	"main/internal/model"
)

type OrderServer interface {
	Connect(ctx context.Context) (<-chan model.Order, error)
	DisConnect(ctx context.Context) error

	SupportOrderType(ctx context.Context) ([]model.OrderType, error)
	PostOrder(ctx context.Context, order model.Order) error
}
