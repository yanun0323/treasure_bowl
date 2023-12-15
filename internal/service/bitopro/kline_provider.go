// documentation
// ohlc: https://github.com/bitoex/bitopro-offical-api-docs/blob/master/api/v3/public/get_ohlc_data.md
package bitopro

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"main/internal/domain"
	"main/internal/entity"
	"main/internal/util"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/yanun0323/decimal"
	"github.com/yanun0323/pkg/logs"
)

const (
	_klineUrl = "/trading-history/"
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
	l               logs.Logger
	ch              chan entity.Kline
	pair            entity.Pair
	targetKlineType entity.KlineType
	cronJob         *cron.Cron
	startAt         int64
}

func NewKlineProvider(ctx context.Context, pair entity.Pair, target entity.KlineType) (domain.KlineProvideServer, error) {
	if _, ok := _klineTypeTrans[target]; !ok {
		return nil, errors.New(fmt.Sprintf("unsupported kline type: %s", target.String()))
	}
	p := &klineProvider{
		l:               logs.Get(ctx).WithField("server", "bitopro kline"),
		pair:            pair,
		targetKlineType: target,
		cronJob:         cron.New(),
		startAt:         time.Now().Unix(),
	}

	p.cronJob.AddFunc(util.CronSpec(), func() {
		p.publishKline(ctx)
	})

	return p, nil
}

func (p *klineProvider) Connect(ctx context.Context, requiredKlineInitCount int) (<-chan entity.Kline, error) {
	result, err := p.requestKline(context.Background(), requiredKlineInitCount)
	if err != nil {
		return nil, errors.Wrap(err, "request kline")
	}

	if len(p.ch) == 0 {
		p.ch = make(chan entity.Kline, len(result)*2)
	}

	for i := range result {
		p.ch <- result[i]
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

	r, err := http.NewRequestWithContext(c, http.MethodGet, p.getApiPath(count), nil)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	ohlc := &OHLC{}
	if err := json.Unmarshal(buf, ohlc); err != nil {
		errors.Wrap(err, "parsing ohlc json data")
	}

	result := make([]entity.Kline, 0, len(ohlc.Data))
	for _, data := range ohlc.Data {
		result = append(result, data.Kline(p.pair, p.targetKlineType))
	}
	return result, nil
}

func (p *klineProvider) publishKline(ctx context.Context) {
	result, err := p.requestKline(ctx, 5)
	if err != nil {
		p.l.WithError(err).Error("request kline")
		return
	}

	for i := range result {
		p.ch <- result[i]
	}
}

func (p *klineProvider) getApiPath(klineCount int) string {
	to := time.Now()
	d := p.targetKlineType.Duration(to.Unix())
	from := to.Add(time.Duration(klineCount) * d)

	return fmt.Sprintf("%s%s%s?resolution=%s&from=%d&to=%d",
		_restfulHost,
		_klineUrl,
		p.pair.Uppercase("_"),
		_klineTypeTrans[p.targetKlineType],
		from.Unix(),
		to.Unix())
}

type OHLC struct {
	Data []OHLCKlineData `json:"data"`
}

type OHLCKlineData struct {
	Timestamp int64           `json:"timestamp"`
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	Volume    decimal.Decimal `json:"volume"`
}

func (d *OHLCKlineData) Kline(p entity.Pair, t entity.KlineType) entity.Kline {
	return entity.Kline{
		Pair:       p,
		Type:       t,
		Source:     entity.KlineSourceBitoPro,
		OpenPrice:  d.Open,
		ClosePrice: d.Close,
		HighPrice:  d.High,
		LowPrice:   d.Low,
		Volume:     d.Volume,
		Timestamp:  d.Timestamp / 1000,
	}
}
