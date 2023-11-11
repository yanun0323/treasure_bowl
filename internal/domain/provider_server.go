package domain

import (
	"context"
	"main/internal/model"
)

type PriceProvideServer interface {
	Connect(ctx context.Context) (<-chan model.Price, error)
	Disconnect(ctx context.Context) error
}

type AssetProvideServer interface {
	Connect(ctx context.Context) (<-chan model.Account, error)
	Disconnect(ctx context.Context) error
}
