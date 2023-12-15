// documentation
// go example: https://github.com/binance/binance-connector-go
// socket: https://github.com/binance/binance-spot-api-docs/blob/master/web-socket-api.md#klines
// socket-stream: https://github.com/binance/binance-spot-api-docs/blob/master/web-socket-streams.md#klinecandlestick-streams
package binance

import (
	"context"
	"main/internal/domain"
	"main/internal/entity"
	"main/internal/util"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"
)

const (
	_klineUrl = "https://api.binance.com"
)

var (
	_klineTypeTrans = map[entity.KlineType]string{
		entity.K1m:  "1m",
		entity.K5m:  "5m",
		entity.K15m: "15m",
		entity.K30m: "30m",
		entity.K1h:  "1h",
		entity.K4h:  "4h",
		entity.K6h:  "6h",
		entity.K12h: "12h",
		entity.K1d:  "1d",
		entity.K1w:  "1w",
		entity.K1M:  "1M",
	}
)

type klineProvider struct {
	l       logs.Logger
	ch      chan entity.Kline
	pair    entity.Pair
	kt      entity.KlineType
	cronJob *cron.Cron
}

func NewKlineProvider(ctx context.Context, pr entity.Pair, target entity.KlineType) (domain.KlineProvideServer, error) {
	if _, ok := _klineTypeTrans[target]; !ok {
		return nil, errors.Errorf("unsupported kline type: %s", target.String())
	}
	p := &klineProvider{
		l:    logs.Get(ctx).WithField("server", "binance kline"),
		pair: pr,
		kt:   target,
	}

	p.cronJob.AddFunc(util.CronSpec(), func() {
		p.publishKline(ctx)
	})

	return p, nil
}

func (p *klineProvider) Connect(ctx context.Context, requiredKlineInitCount int) (<-chan entity.Kline, error) {
	kls, err := p.requestKline(ctx, requiredKlineInitCount)
	if err != nil {
		return nil, errors.Wrap(err, "request kline")
	}

	if len(p.ch) == 0 {
		p.ch = make(chan entity.Kline, len(kls)*2)
	}

	for i := range kls {
		p.ch <- kls[i]
	}

	p.cronJob.Start()
	return p.ch, nil
}

func (p *klineProvider) Disconnect(ctx context.Context) error {
	p.cronJob.Stop()
	return nil
}

func (p *klineProvider) requestKline(ctx context.Context, count int) ([]entity.Kline, error) {
	c, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	client, err := NewClient(_klineUrl)
	if err != nil {
		return nil, errors.Wrap(err, "new client")
	}

	klines, err := client.NewKlinesService().Symbol(p.pair.Uppercase()).Interval(_klineTypeTrans[p.kt]).Do(c)
	if err != nil {
		return nil, errors.Errorf("new klines service, err: %+v", err)
	}

	result := make([]entity.Kline, 0, len(klines))
	for _, kl := range klines {
		k, err := p.toEntityKline(kl)
		if err != nil {
			return nil, errors.Wrap(err, "transfer kline data")
		}
		result = append(result, k)
	}
	return result, nil
}

func (p *klineProvider) publishKline(ctx context.Context) {
	kls, err := p.requestKline(ctx, 5)
	if err != nil {
		p.l.WithError(err).Error("request kline")
		return
	}

	for i := range kls {
		p.ch <- kls[i]
	}
}

func (p *klineProvider) toEntityKline(resp *binance_connector.KlinesResponse) (entity.Kline, error) {
	open, err := decimal.New(resp.Open)
	if err != nil {
		return entity.Kline{}, errors.Errorf("new decimal open, err: %+v", err)
	}

	close, err := decimal.New(resp.Close)
	if err != nil {
		return entity.Kline{}, errors.Errorf("new decimal close, err: %+v", err)
	}

	high, err := decimal.New(resp.High)
	if err != nil {
		return entity.Kline{}, errors.Errorf("new decimal high, err: %+v", err)
	}

	low, err := decimal.New(resp.Low)
	if err != nil {
		return entity.Kline{}, errors.Errorf("new decimal low, err: %+v", err)
	}

	return entity.Kline{
		Pair:       p.pair,
		Type:       p.kt,
		Source:     entity.KlineSourceBinance,
		OpenPrice:  open,
		ClosePrice: close,
		HighPrice:  high,
		LowPrice:   low,
		Timestamp:  int64(resp.CloseTime),
	}, nil
}
