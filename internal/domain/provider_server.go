package domain

import (
	"context"
	"main/internal/model"
)

type KlineProvideServer interface {
	Connect(ctx context.Context, types ...model.KlineType) (<-chan model.Kline, error)
	Disconnect(ctx context.Context) error
}

type AssetProvideServer interface {
	Connect(ctx context.Context) (<-chan model.Account, error)
	Disconnect(ctx context.Context) error
}
