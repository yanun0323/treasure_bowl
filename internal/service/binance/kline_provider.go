// documentation
// go example: https://github.com/binance/binance-connector-go
// socket: https://github.com/binance/binance-spot-api-docs/blob/master/web-socket-api.md#klines
// socket-stream: https://github.com/binance/binance-spot-api-docs/blob/master/web-socket-streams.md#klinecandlestick-streams
package binance

import (
	"context"
	"main/internal/domain"
	"main/internal/entity"

	"github.com/binance/binance-connector-go"
	"github.com/pkg/errors"
	"github.com/yanun0323/pkg/logs"
)

type klineProvider struct {
	l    logs.Logger
	ch   chan entity.Kline
	pair entity.Pair
}

func NewKlineProvider(ctx context.Context, pr entity.Pair) (domain.KlineProvideServer, error) {
	return &klineProvider{
		l:    logs.Get(ctx).WithField("server", "binance kline"),
		pair: pr,
	}, nil
}

func (p *klineProvider) Connect(ctx context.Context, requiredKlineInitCount int) (<-chan entity.Kline, error) {
	baseURL := "https://api.binance.com"

	client := binance_connector.NewClient("", "", baseURL)

	// TODO: Implement channel

	// Klines
	klines, err := client.NewKlinesService().Symbol(p.pair.Uppercase()).Interval("1m").Do(ctx)
	if err != nil {
		return nil, errors.Errorf("new klines service, err: %+v", err)
	}
	p.l.Debug(binance_connector.PrettyPrint(klines))
	return nil, nil
}

func (p *klineProvider) Disconnect(ctx context.Context) error {
	return nil
}
