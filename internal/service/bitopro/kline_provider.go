package bitopro

import (
	"context"

	"main/internal/domain"
	"main/internal/model"
)

type KlineProvider struct {
}

func NewKlineProvider(pair model.Pair) domain.KlineProvideServer {

	return &KlineProvider{}
}

func (p *KlineProvider) Connect(ctx context.Context, types ...model.KlineType) (<-chan model.Kline, error) {
	return make(<-chan model.Kline), nil
}

func (p *KlineProvider) Disconnect(ctx context.Context) error {
	return nil
}
