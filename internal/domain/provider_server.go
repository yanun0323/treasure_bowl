package domain

import (
	"context"
	"main/internal/model"
)

type KLineProvideServer interface {
	Connect(ctx context.Context) (<-chan model.KLine, error)
	Disconnect(ctx context.Context) error
}

type AssetProvideServer interface {
	Connect(ctx context.Context) (<-chan model.Account, error)
	Disconnect(ctx context.Context) error
}
