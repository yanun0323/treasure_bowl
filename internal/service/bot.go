package service

import (
	"context"
	"main/internal/domain"
	"main/internal/model"

	"github.com/yanun0323/pkg/logs"
)

type bot struct {
	Pair                string
	Log                 logs.Logger
	OrderServer         domain.OrderServer
	PriceProviderServer domain.PriceProvideServer
	AssetProviderServer domain.AssetProvideServer
	StrategyServer      domain.StrategyServer
}

func NewBot() (bot, error) {
	// TODO: Implement me
	return bot{}, nil
}

func (b *bot) Run(ctx context.Context) error {
	if err := b.setup(ctx); err != nil {
		return err
	}

	priceCh, err := b.PriceProviderServer.Connect(ctx)
	if err != nil {
		return err
	}

	assetCh, err := b.AssetProviderServer.Connect(ctx)
	if err != nil {
		return err
	}

	signal, err := b.StrategyServer.Connect(ctx)
	if err != nil {
		return err
	}

	go consumePrice(ctx, priceCh, b.StrategyServer)
	go consumeAsset(ctx, assetCh, b.StrategyServer)
	go consumeSignal(ctx, signal, b.OrderServer, b.StrategyServer)

	return nil
}

func (b *bot) Shutdown(ctx context.Context) error {
	if err := b.StrategyServer.Disconnect(ctx); err != nil {
		return err
	}

	if err := b.PriceProviderServer.Disconnect(ctx); err != nil {
		return err
	}

	if err := b.AssetProviderServer.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (b *bot) setup(ctx context.Context) error {
	orders, err := b.OrderServer.Orders(ctx)
	if err != nil {
		return err
	}
	b.StrategyServer.PushOrders(ctx, orders...)

	types, err := b.OrderServer.SupportOrderType(ctx)
	if err != nil {
		return err
	}
	b.StrategyServer.PushSupportedOrderTypes(ctx, types...)

	return nil
}

func consumePrice(ctx context.Context, ch <-chan model.Price, strategy domain.StrategyServer) {
	for {
		select {
		case price := <-ch:
			strategy.PushPrices(ctx, price)
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

func consumeSignal(ctx context.Context, signal <-chan model.Order, orderServer domain.OrderServer, strategy domain.StrategyServer) {
	for {
		select {
		case order := <-signal:
			orders, err := orderServer.PostOrder(ctx, order)
			if err != nil {
				logs.Get(ctx).Errorf("set order '%s', err: %s", order.ID, err.Error())
				continue
			}
			strategy.PushOrders(ctx, orders...)
		case <-ctx.Done():
			return
		}
	}
}
