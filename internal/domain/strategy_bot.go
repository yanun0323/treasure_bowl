package domain

import (
	"context"
)

type StrategyBot interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
