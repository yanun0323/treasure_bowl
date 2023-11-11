package domain

import (
	"context"
	"main/internal/model"
)

type StrategyServer interface {
	PushPrice(ctx context.Context, pair string, price model.Price) error
}
