package domain

import (
	"context"
	"time"

	"main/internal/entity"
)

type KlineProvideServer interface {
	Connect(ctx context.Context, requiredKlineInitCount int) (<-chan entity.Kline, error)
	Disconnect(ctx context.Context) error
	History(ctx context.Context, start, end time.Time) (<-chan entity.Kline, error)
}
