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
	KLineProviderServer domain.KLineProvideServer
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

	kLineCh, err := b.KLineProviderServer.Connect(ctx)
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

	go consumeKLine(ctx, kLineCh, b.StrategyServer)
	go consumeAsset(ctx, assetCh, b.StrategyServer)
	go consumeOrder(ctx, orderCh, b.StrategyServer)
	go consumeSignal(ctx, signal, b.OrderServer)

	return nil
}

func (b *bot) Shutdown(ctx context.Context) error {
	if err := b.StrategyServer.Disconnect(ctx); err != nil {
		return err
	}

	if err := b.KLineProviderServer.Disconnect(ctx); err != nil {
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
	b.StrategyServer.PushSupportedOrderTypes(ctx, types...)

	return nil
}

func consumeKLine(ctx context.Context, ch <-chan model.KLine, strategy domain.StrategyServer) {
	for {
		select {
		case kLine := <-ch:
			strategy.PushKLines(ctx, kLine)
		case <-ctx.Done():
			return
		}
	}
}

func consumeAsset(ctx context.Context, ch <-chan model.Account, strategy domain.StrategyServer) {
	for {
		select {
		case account := <-ch:
			strategy.PushAssets(ctx, account)
		case <-ctx.Done():
			return
		}
	}
}

func consumeOrder(ctx context.Context, ch <-chan model.Order, strategy domain.StrategyServer) {
	for {
		select {
		case order := <-ch:
			strategy.PushOrders(ctx, order)
		case <-ctx.Done():
			return
		}
	}
}

func consumeSignal(ctx context.Context, signal <-chan model.Order, orderServer domain.OrderServer) {
	for {
		select {
		case order := <-signal:
			if err := orderServer.PostOrder(ctx, order); err != nil {
				logs.Get(ctx).Errorf("set order '%s', err: %s", order.ID, err.Error())
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}