package domain

import (
	"context"

	"main/internal/entity"
)

type KlineProvideServer interface {
	Connect(ctx context.Context, requiredKlineInitCount int) (<-chan entity.Kline, error)
	Disconnect(ctx context.Context) error
}
