package bitopro

import (
	"context"
	"strings"

	"main/internal/domain"
	"main/internal/model"
	"main/internal/util"

	"github.com/bitoex/bitopro-api-go/pkg/ws"
	"github.com/pkg/errors"
	"github.com/yanun0323/pkg/logs"
)

type KlineProvider struct {
	l              logs.Logger
	tickers        chan model.Kline
	close          chan struct{}
	cancel         context.CancelFunc
	pair           model.Pair
	supportedTypes []model.KlineType
}

func NewKlineProvider(pair model.Pair) domain.KlineProvideServer {
	return &KlineProvider{
		l:       logs.New("bitopro kline provider", util.LogLevel()),
		tickers: make(chan model.Kline, 50),
	}
}

func (p *KlineProvider) Connect(ctx context.Context, types ...model.KlineType) (<-chan model.Kline, error) {
	ticker, close := ws.NewPublicWs().RunTickerWsConsumer(ctx, []string{p.pair.String("_")})
	if ticker == nil {
		return nil, errors.New("connect to ticker")
	}
	p.close = close

	p.supportedTypes = types

	c, cancel := context.WithCancel(ctx)
	go p.consumeTicker(c, ticker)
	p.cancel = cancel
	return p.tickers, nil
}

func (p *KlineProvider) Disconnect(ctx context.Context) error {
	close(p.close)
	p.cancel()
	return nil
}

func (p *KlineProvider) consumeTicker(ctx context.Context, ticker <-chan ws.TickerData) {
	for {
		select {
		case t := <-ticker:
			kline, err := parseTickerDataToKline(t)
			if err != nil {
				p.l.WithError(err).Error("consume ticker")
				continue
			}
			p.tickers <- kline
		case <-ctx.Done():
			return
		}
	}
}

func parseTickerDataToKline(td ws.TickerData) (model.Kline, error) {
	if td.Err != nil {
		return model.Kline{}, errors.Wrap(td.Err, "ticker data error")
	}

	pairs := strings.Split(td.Pair, "_")

	return model.Kline{
		Pair:      model.NewPair(pairs[0], pairs[1]),
		Timestamp: td.Timestamp,
	}, nil
}

/* {
   "event": "TICKER",
   "pair": "BTC_TWD",
   "lastPrice": "1",
   "isBuyer": true,
   "priceChange24hr": "1",
   "volume24hr": "1",
   "high24hr": "1",
   "low24hr": "1",
   "timestamp": 1136185445000,
   "datetime": "2006-01-02T15:04:05.700Z"
 } */
