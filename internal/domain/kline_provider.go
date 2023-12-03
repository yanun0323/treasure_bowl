package domain

import (
	"context"

	"main/internal/model"
)

type KlineProvideServer interface {
	Connect(ctx context.Context, requiredKlineInitCount int) (<-chan model.Kline, error)
	Disconnect(ctx context.Context) error
}
