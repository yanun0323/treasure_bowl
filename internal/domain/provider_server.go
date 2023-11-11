package domain

import (
	"context"
	"main/internal/model"
)

type PriceProviderServer interface {
	SubscribePrice(ctx context.Context, pair string) (<-chan model.Price, error)
}

type AssetProviderServer interface {
	SubscribeAsset(ctx context.Context) (<-chan model.Account, error)
}
