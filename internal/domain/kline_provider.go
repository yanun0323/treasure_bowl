package domain

import (
	"context"

	"main/internal/entity"
)

type KlineProvideServer interface {
	Connect(ctx context.Context, spans ...entity.Span) (<-chan entity.Kline, error)
	Disconnect(ctx context.Context) error
}
