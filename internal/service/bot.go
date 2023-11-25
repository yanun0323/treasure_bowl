package service

import (
	"context"

	"main/internal/domain"
	"main/internal/model"

	"github.com/yanun0323/pkg/logs"
)

type bot struct {
	ID                  string
	Pair                string
	Log                 logs.Logger
	OrderServer         domain.OrderServer
	KlineProviderServer domain.KlineProvideServer
	AssetProviderServer domain.AssetProvideServer
	StrategyServer      domain.StrategyServer
}

func NewBot(id string, pair string) (bot, error) {
	// TODO: Implement me
	return bot{
		ID:   id,
		Pair: pair,
	}, nil
}

func (b *bot) Run(ctx context.Context) error {
	if err := b.setup(ctx); err != nil {
		return err
	}

	klineCh, err := b.KlineProviderServer.Connect(ctx)
	if err != nil {
		return err
	}

	assetCh, err := b.AssetProviderServer.Connect(ctx)
	if err != nil {
		return err
	}

	orderCh, err := b.OrderServer.Connect(ctx)
	if err != nil {
		return err
	}

	signal, err := b.StrategyServer.Connect(ctx)
	if err != nil {
		return err
	}

	go consumeKline(ctx, klineCh, b.StrategyServer)
	go consumeAsset(ctx, assetCh, b.StrategyServer)
	go consumeOrder(ctx, orderCh, b.StrategyServer)
	go consumeSignal(ctx, signal, b.OrderServer)

	return nil
}

func (b *bot) Shutdown(ctx context.Context) error {
	if err := b.StrategyServer.Disconnect(ctx); err != nil {
		return err
	}

	if err := b.KlineProviderServer.Disconnect(ctx); err != nil {
		return err
	}

	if err := b.AssetProviderServer.Disconnect(ctx); err != nil {
		return err
	}

	if err := b.OrderServer.DisConnect(ctx); err != nil {
		return err
	}

	return nil
}

func (b *bot) setup(ctx context.Context) error {
	types, err := b.OrderServer.SupportOrderType(ctx)
	if err != nil {
		return err
	}
	b.StrategyServer.PushSupportedOrderType(ctx, types...)

	return nil
}

func consumeKline(ctx context.Context, ch <-chan model.Kline, strategy domain.StrategyServer) {
	for {
		select {
		case kline := <-ch:
			strategy.PushKline(ctx, kline)
		case <-ctx.Done():
			return
		}
	}
}

func consumeAsset(ctx context.Context, ch <-chan model.Account, strategy domain.StrategyServer) {
	for {
		select {
		case account := <-ch:
			strategy.PushAsset(ctx, account)
		case <-ctx.Done():
			return
		}
	}
}

func consumeOrder(ctx context.Context, ch <-chan model.Order, strategy domain.StrategyServer) {
	for {
		select {
		case order := <-ch:
			strategy.PushOrder(ctx, order)
		case <-ctx.Done():
			return
		}
	}
}

func consumeSignal(ctx context.Context, signal <-chan model.Order, orderServer domain.OrderServer) {
	for {
		select {
		case order := <-signal:
			if err := orderServer.PushOrder(ctx, order); err != nil {
				logs.Get(ctx).Errorf("set order '%s', err: %s", order.ID, err.Error())
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}
