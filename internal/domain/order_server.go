package domain

import (
	"context"
	"main/internal/model"
)

type OrderServer interface {
	SupportOrderType(ctx context.Context) ([]model.OrderType, error)
	Orders(ctx context.Context) ([]model.Order, error)
	PostOrder(ctx context.Context, order model.Order) ([]model.Order, error)
}
