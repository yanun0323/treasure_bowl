package domain

import (
	"context"
	"github.com/labstack/echo/v4"
)

type StrategyBot interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error

	GetInfo(c echo.Context) error
}
